# LEFFINDO Resi Scanner

A high performance logistics and shipment tracking tool built with Go and Fyne, developed specially for PT LEFFINDO JAYA MANDIRI based in Batam. This application is designed to manage shipment data via Excel imports and provides a high-speed scanning interface optimized for Android tablets and physical barcode scanners.

## 🚀 Features

- **High-Speed Scanning:** Features a dedicated input layer optimized for industrial hardware scanners, ensuring rapid and reliable data capture on Android devices.
- **SQLite Backend:** Persistent storage with a relational schema for Shipments and Locations.
- **Excel Integration:** Import bulk shipment data directly from Excel files (with a proper format).
- **Smart UI:** Auto-scrolling, row highlighting, and automatic focus management for continuous scanning workflows.
- **Cross-Platform:** Currently optimized for Android (Production).

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **UI Framework:** [Fyne v2](https://fyne.io)
- **Database:** SQLite3 (via `github.com/mattn/go-sqlite3`)
- **Build Tool:** `fyne-cross` for Android compilation.

### Prerequisites
- Go 1.21+
- Android NDK (for mobile builds)
- `fyne-cross`

## 📥 Installation & Setup Guide

### 1. Tablet Preparation (Android)
To install the application directly via APK, you must allow your tablet to install apps from outside the Google Play Store:
* Go to **Settings > Apps > Special app access**.
* Select **Install unknown apps**.
* Toggle **Allow** for your Browser (e.g., Chrome) or File Manager.

### 2. Physical Scanner Setup
For the best high-speed scanning experience on the tablet:
* **Connect the Scanner:** Plug in your USB scanner or pair your Bluetooth scanner.
* **Configure Keyboard:** Go to **Settings > Additional Settings > Languages & Input > Physical Keyboard**.
* **Disable Virtual Keyboard:** Turn **OFF** "Show virtual keyboard." This prevents the on-screen keyboard from blocking the UI during scanning.

### 3. Application Launch
1. Download the latest `.apk` from the **Releases** section of this repository.
2. Open the file and tap **Install**.
3. Upon first launch, the app will automatically initialize the `resi-scanner.db` database in the internal storage.
4. Use the **Excel Import** button to load your initial shipment data (ensure the file follows the required template format).

## 🔄 How to Update
When a new version is released:
1. Go to the **Releases** page on GitHub.
2. Download the new `.apk`.
3. Install it over the existing app. Your database (`.db` file) will remain safe and will not be deleted during the update.

---
*Developed for PT LEFFINDO JAYA MANDIRI - Batam, Indonesia.*
