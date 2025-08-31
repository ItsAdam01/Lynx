---
title: "Isolated Lab Testing (Experimental)"
weight: 4
---

# Experimental: Isolated Lab Setup

If you want to test Lynx FIM without modifying your system directories or adding the binary to your global `$PATH`, you can set up an isolated "Lab" directory. This is exactly how I did my final testing for this project.

## 🧪 The "Single Directory" Setup

Follow these steps to create a self-contained testing environment.

### 1. Create the Lab Directory
```bash
mkdir lynx-lab && cd lynx-lab
```

### 2. Prepare the Binary
Build Lynx from the root of the repository and copy it into your lab:
```bash
# From the project root
make build
cp bin/lynx ./lynx-lab/
cd lynx-lab
```

### 3. Create Dummy Files to Watch
```bash
mkdir watched_files
echo "initial content" > watched_files/test.txt
```

### 4. Initialize and Configure Locally
Run the init command inside the lab:
```bash
./lynx init
```

Open `config.yaml` and change the `paths_to_watch` to point only to your lab directory:
```yaml
paths_to_watch:
  - "./watched_files"
```

### 5. Run the Detection Cycle
In your terminal, set the secret and create the baseline:
```bash
export LYNX_HMAC_SECRET="lab-testing-secret"
./lynx baseline -o lab_baseline.json
```

Now, start the monitor:
```bash
./lynx start -b lab_baseline.json
```

### 6. Test the Detection
Open a **second terminal window**, navigate to the lab directory, and trigger an event:
```bash
echo "tampering with file" >> watched_files/test.txt
```

Go back to your first terminal. You should see the `CRITICAL: File modified` alert appear instantly in both your console and the `lynx.log` file.

## 💡 Why Test This Way?
- **Safety:** You aren't touching sensitive system files like `/etc/passwd`.
- **Speed:** It's much faster to iterate on configuration changes.
- **Portability:** You can delete the `lynx-lab` folder when you're done, leaving your system completely clean.
