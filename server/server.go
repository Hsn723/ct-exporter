package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Hsn723/certspotter-client/api"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "certificate_transparency"
)

var (
	issuanceGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, "", "issuance"),
			Help: "issuance observed in certificate transparency log",
		},
		[]string{"serial_no", "issuer_cn", "dns_names", "not_after", "type", "target"},
	)
)

type CTExporter struct {
	Addr        string
	Client      api.CertspotterClient
	Env         *well.Environment
	PositionDir string
}

func (s *CTExporter) createLabel(issuance api.Issuance) ([]string, error) {
	dnsnames := ""
	if len(issuance.Domains) > 0 {
		dnsnames = strings.Join(issuance.Domains, ",")
	}
	notAfter, err := time.Parse(time.RFC3339, issuance.NotAfter)
	if err != nil {
		return nil, err
	}

	return []string{
		strconv.FormatUint(issuance.ID, 10),
		issuance.Issuer.Name,
		dnsnames,
		strconv.FormatInt(notAfter.Unix(), 10),
		issuance.Cert.Type,
	}, nil
}

func (s *CTExporter) getPosition(filename string) (pos uint64, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, nil
	}
	return strconv.ParseUint(strings.TrimSpace(string(content)), 10, 64)
}

func (s *CTExporter) writePosition(filename string, position uint64) error {
	tmpFile, err := ioutil.TempFile(filepath.Dir(filename), "")
	if err != nil {
		return err
	}
	_, err = tmpFile.WriteString(strconv.FormatUint(position, 10))
	if err != nil {
		return err
	}
	fi, err := tmpFile.Stat()
	if err != nil {
		return err
	}
	if fi.Size() == 0 {
		return nil
	}
	return os.Rename(tmpFile.Name(), filename)
}

func (s *CTExporter) recordIssuances(issuances []api.Issuance, target string) {
	for _, issuance := range issuances {
		labels, err := s.createLabel(issuance)
		if err != nil {
			_ = log.Error("error creating labels", map[string]interface{}{
				"target": target,
				"error":  err.Error(),
			})
		}
		labels = append(labels, target)
		notBefore, err := time.Parse(time.RFC3339, issuance.NotBefore)
		if err != nil {
			_ = log.Error("error parsing notBefore date", map[string]interface{}{
				"target": target,
				"error":  err.Error(),
			})
		}
		issuanceGauge.WithLabelValues(labels...).Set(float64(notBefore.Unix()))
	}
}

func (s *CTExporter) handleProbe(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := filepath.Join(s.PositionDir, target)
	pos, err := s.getPosition(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if pos > 0 {
		_ = log.Info("using position", map[string]interface{}{
			"target":   target,
			"position": pos,
		})
	}

	issuances, err := s.Client.GetIssuances(target, true, true, pos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(issuanceGauge)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	s.recordIssuances(issuances, target)
	h.ServeHTTP(w, r)
	if len(issuances) == 0 {
		_ = log.Info("no new issuances observed", map[string]interface{}{
			"target": target,
		})
		return
	}

	if err := s.writePosition(filename, issuances[len(issuances)-1].ID); err != nil {
		_ = log.Error("error saving positions", map[string]interface{}{
			"target": target,
			"error":  err.Error(),
		})
	}
}

func (s *CTExporter) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/probe", s.handleProbe).Methods(http.MethodGet)
	server := &well.HTTPServer{
		Server: &http.Server{
			Handler: router,
			Addr:    s.Addr,
		},
		AccessLog: log.DefaultLogger(),
		Env:       s.Env,
	}

	_ = log.Info("started listening", map[string]interface{}{
		"address": server.Server.Addr,
	})
	return server.ListenAndServe()
}
