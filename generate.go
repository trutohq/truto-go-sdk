package truto

// Regenerate the typed models from spec/openapi.json, with the fixes in
// spec/overlay.json applied:
//
//	go generate ./...
//
//go:generate sh -c "python3 codegen/normalize_spec.py spec/openapi.json spec/overlay.json > codegen/models.json && python3 codegen/generate_go.py codegen/models.json > models_gen.go && gofmt -w models_gen.go"
