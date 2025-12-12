.PHONY: proto proto-lint proto-breaking clean-proto test test-backend test-compute test-coverage test-backend-coverage test-compute-coverage test-backend-quick test-compute-quick coverage-summary

proto:
	@echo "Generating protobuf code..."
	buf generate
	@if [ -f compute-client/src/proto/compute/v1/service_pb2_grpc.py ]; then \
		sed -i 's/from compute\.v1 import/from proto.compute.v1 import/g' compute-client/src/proto/compute/v1/service_pb2_grpc.py; \
	fi
	@echo "Proto generation complete!"

proto-lint:
	@echo "Linting proto files..."
	buf lint

proto-breaking:
	@echo "Checking for breaking changes..."
	buf breaking --against '.git#branch=main'

clean-proto:
	@echo "Cleaning generated proto code..."
	rm -rf backend/pkg/proto/compute
	rm -rf compute-client/src/proto/*_pb2.py
	rm -rf compute-client/src/proto/*_pb2_grpc.py
	@echo "Clean complete!"

# ==================== Testing ====================

test: test-backend test-compute
	@echo "All tests complete!"

test-backend:
	@echo "Running backend tests..."
	cd backend && go test -v ./...
	@echo "Backend tests complete!"

test-compute:
	@echo "Running compute-client tests..."
	cd compute-client && python -m pytest tests/ -v
	@echo "Compute-client tests complete!"

test-coverage: test-backend-coverage test-compute-coverage
	@echo "All coverage reports generated!"

test-backend-coverage:
	@echo "Running backend tests with coverage..."
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Backend coverage report: backend/coverage.html"

test-compute-coverage:
	@echo "Running compute-client tests with coverage..."
	cd compute-client && python -m pytest tests/ --cov=src --cov-report=html --cov-report=term
	@echo "Compute-client coverage report: compute-client/htmlcov/index.html"

test-backend-quick:
	@cd backend && go test ./...

test-compute-quick:
	@cd compute-client && python -m pytest tests/

coverage-summary:
	@cd backend && go test -cover ./...