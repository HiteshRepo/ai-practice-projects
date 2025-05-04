.PHONY: install-portaudio
install-portaudio:
	@echo "Installing PortAudio..."
	@brew install pkg-config
	@brew install portaudio

.PHONY: find-portaudio-file
find-portaudio-file:
	@echo "Finding the PortAudio library file..."
	@find /usr/local/lib /opt/homebrew/lib -name "portaudio-2.0.pc"

.PHONY: win-install-portaudio
win-install-portaudio:
	@echo "Installing PortAudio on Windows..."
	@echo "Downloading and installing pkg-config..."
	@powershell -Command "if (-not (Get-Command -ErrorAction SilentlyContinue choco)) { Write-Host 'Installing Chocolatey...'; Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1')) }"
	@powershell -Command "choco install -y pkgconfiglite"
	@echo "Downloading and installing PortAudio..."
	@powershell -Command "Invoke-WebRequest -Uri 'http://files.portaudio.com/archives/pa_stable_v190700_20210406.tgz' -OutFile 'portaudio.tgz'"
	@powershell -Command "tar -xzf portaudio.tgz"
	@cd portaudio && cmake -G "MinGW Makefiles" . && mingw32-make && mingw32-make install
	@echo "PortAudio installation complete!"

.PHONY: win-find-portaudio-file
win-find-portaudio-file:
	@echo "Finding the PortAudio library file..."
	@powershell -Command "Get-ChildItem -Path 'C:\Program Files', 'C:\Program Files (x86)', 'C:\MinGW', 'C:\msys64' -Filter 'portaudio-2.0.pc' -Recurse -ErrorAction SilentlyContinue | ForEach-Object { $_.FullName }"