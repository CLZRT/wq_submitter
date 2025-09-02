# Makefile for wq_submitter project

# --- Variables ---
APP_NAME := wq_submitter
BINARY := $(APP_NAME)

# Install locations (sudo required)
INSTALL_DIR := /usr/local/$(APP_NAME)
CONFIG_DIR := $(INSTALL_DIR)/configs
CONFIG_SRC := ./configs/config.yaml
SYSTEMD_FILE := /etc/systemd/system/$(APP_NAME).service

# --- Main Targets ---
.PHONY: setup-env start stop restart status


help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Development Targets:"
	@echo "  setup-env    Generate environment variables. Usage: eval \"$$(make setup-env)\""
	@echo ""
	@echo "Deployment Targets (requires sudo):"
	@echo "  install      Install the application as a systemd service."
	@echo "  uninstall    Remove the application and systemd service."
	@echo "  start        Start the service."
	@echo "  stop         Stop the service."
	@echo "  restart      Restart the service."
	@echo "  status       Show the service status."

setup-env:
	@# This target outputs a shell script to be used with eval.
	@echo "echo '--- Setting up environment variables ---';" \
		"TOKEN=$$(openssl rand -base64 32 | tr -d '\n');" \
		"export WQS_SECRET_TOKEN='$$TOKEN';" \
		"echo '✔ Secret token (WQS_SECRET_TOKEN) generated.';" \
		"read -p 'Enter Brain Username: ' BRAIN_USERNAME;" \
		"export BRAIN_USERNAME;" \
		"read -sp 'Enter Brain Password: ' BRAIN_PASSWORD;" \
		"export BRAIN_PASSWORD;" \
		"echo;" \
		"echo '✔ Brain credentials set for this session.';"


install:
	@echo "--> Installing application (sudo required)..."
	@sudo mkdir -p $(INSTALL_DIR)
	@sudo mkdir -p $(CONFIG_DIR)
	@sudo cp $(BINARY) $(INSTALL_DIR)/
	@sudo cp $(CONFIG_SRC) $(CONFIG_DIR)/
	@echo "✔ Binary and config copied to $(INSTALL_DIR)"
	@echo "--> Creating systemd service..."
	@printf '[Unit]\nDescription=$(APP_NAME) Service\nAfter=network.target\n\n[Service]\nType=simple\nUser=root\nGroup=root\nWorkingDirectory=$(INSTALL_DIR)\nExecStart=$(INSTALL_DIR)/$(BINARY)\nRestart=on-failure\nRestartSec=5s\n\n[Install]\nWantedBy=multi-user.target\n' | sudo tee $(SYSTEMD_FILE) > /dev/null
	@echo "✔ Systemd service file created: $(SYSTEMD_FILE)"
	@sudo systemctl daemon-reload
	@sudo systemctl enable $(APP_NAME)
	@sudo systemctl start $(APP_NAME)
	@echo "✔ Service installed and started. Use 'make status' to check."

uninstall:
	@echo "--> Uninstalling application (sudo required)..."
	@sudo systemctl stop $(APP_NAME) || true
	@sudo systemctl disable $(APP_NAME) || true
	@echo "✔ Service stopped and disabled."
	@sudo rm -f $(SYSTEMD_FILE)
	@sudo systemctl daemon-reload
	@echo "✔ Systemd service file removed."
	@sudo rm -rf $(INSTALL_DIR)
	@echo "✔ Application directory removed."
	@echo "✔ Uninstallation complete."

start:
	@sudo systemctl start $(APP_NAME)
	@echo "Service started."

stop:
	@sudo systemctl stop $(APP_NAME)
	@echo "Service stopped."

restart:
	@sudo systemctl restart $(APP_NAME)
	@echo "Service restarted."

status:
	@sudo systemctl status $(APP_NAME)