# Lynx FIM: A Cybersecurity Learning Project

Lynx FIM is a host-based intrusion detection agent (HIDS) I built to understand the fundamentals of file integrity monitoring and real-time system alerting in Go. 

This repository contains the full source code and a detailed documentation site covering my learning journey from June to August 2025.

---

## 📖 Documentation Site

I have created a comprehensive documentation site using Hugo. You can access it locally:
1. Navigate to `site/` and run `hugo server -D`.
2. Visit **`http://localhost:1313/`**.

### Browse Docs on GitHub:
- **[🚀 Usage Guide](site/content/docs/usage/installation.md)**
  - [Installation & Setup](site/content/docs/usage/installation.md)
  - [Command Reference](site/content/docs/usage/commands.md)
  - [General Features](site/content/docs/usage/features.md)
  - [Isolated Lab Testing](site/content/docs/usage/isolated_testing.md)
- **[💻 Development & Research](site/content/docs/development/technical_specs.md)**
  - [Technical Specifications](site/content/docs/development/technical_specs.md)
  - [Implementation Story](site/content/docs/development/implementation_story.md)
  - [Performance Analysis](site/content/docs/development/performance.md)
  - [Proof of Concept](site/content/docs/development/demonstration.md)

---

## 🧪 Quick Test: The Isolated Lab

If you want to see Lynx FIM in action without affecting your system, follow this temporary process:

```bash
# 1. Prepare Workspace
mkdir -p /tmp/lynx-lab && cd /tmp/lynx-lab
# (From project root)
make build && cp bin/lynx /tmp/lynx-lab/
cd /tmp/lynx-lab

# 2. Create Dummy Data
mkdir watched_dirs && echo "secret info" > watched_dirs/top_secret.txt

# 3. Initialize and Configure
./lynx init
sed -i 's|/etc/ssh|./watched_dirs|g' config.yaml

# 4. Set Secret and Baseline
export LYNX_HMAC_SECRET="lab-secret-123"
./lynx baseline -o lab_baseline.json

# 5. Start Monitoring (blocks terminal)
./lynx start -b lab_baseline.json
```

**In a second terminal:**
```bash
echo "tampered!" >> /tmp/lynx-lab/watched_dirs/top_secret.txt
```
*You will see the critical alert immediately in Terminal 1.*

---

- **Identity:** Adam Atienza
- **Timeline:** June 2025 – August 2025
- **Goal:** Learn by building a professional-grade security tool from scratch.
