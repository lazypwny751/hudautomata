all: backend-build frontend-build firmware-build
clean: backend-clean frontend-clean firmware-clean

backend-build:
	@echo "Backend build."

backend-clean:
	@echo "Backend clean."

frontend-build:
	@echo "Frontend build."

frontend-clean:
	@echo "Frontend clean."

firmware-build:
	@echo "Firmware build."

firmware-clean:
	@echo "Firmware clean."

.PHONY: backend-build frontend-build firmware-build
