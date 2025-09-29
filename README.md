# OcelotMDM

![wip](https://img.shields.io/badge/status-WIP-yellow)

*OcelotMDM* it's a lightweight, modular Mobile Device Management system designed embedded devices with limited resources, deployed in challenging environments such as public transport fleets (e.g. buses) with unstable or intermittent connectivity.

---

## 📌 Key Features

- ⚡ **Minimal Footprint**: Designed for embedded Linux and Android devices with low resources.
- 🔒 **Secure Enrollment**: Enroll new devices via OTP, secure certificate exchange, and VPN tunnel setup.
- 🌐 **Distributed Architecture**: Supports operation over unstable networks with persistent command queues.
- 🗂️ **Modular Design**: Separate backend, agent, and dashboard for easy maintenance and scaling.
- 📡 **Remote Command Execution**: Agents execute remote commands, update configurations, and handle OTA updates.
- 📊 **Centralized Logging**: Uses Elasticsearch for device logs, telemetry, and audit trails.
- 🕵️‍♂️ **Fleet Monitoring Dashboard**: Web dashboard to manage devices, monitor health, and push updates.

---

## 🚦 Roadmap (Draft)

- [x] Backend core API
- [x] MQTT broker setup
- [x] Linux agent MVP
- [ ] Android agent MVP
- [x] Enrollment flow (OTP, cert exchange)
- [x] VPN tunnel automation
- [ ] Web dashboard basic version
- [ ] Logging and metrics pipeline (Elastic)
- [x] Deployment scripts & NGINX config

---

## 🔐 Security Model

- Each device enrolls with a one-time token (OTP) for initial trust.
- Server provisions certificates for VPN and secure MQTT communication.
- Commands and telemetry encrypted in transit (TLS).

---

## 📚 Project Context

*OcelotMDM* was originally developed as part of a university thesis project focused on embedded fleet management for public transport vehicles.  
The goal is to provide a robust, scalable, and open solution where existing MDMs fail due to network instability, limited hardware, or overly complex architectures.

---

## ✨ Credits

Maintained by **Matteo Brusarosco**  
Initial development for thesis at **University of Trento**

---
