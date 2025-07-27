# Tyro: Interactive DICOM Administrative Tool
<!-- 
![Go Build](https://github.com/yourusername/tyro/workflows/Go/badge.svg) <!-- Placeholder: You'll add your CI/CD badge here later -->
<!--![Latest Release](https://img.shields.io/github/v/release/yourusername/tyro) Placeholder -->
![License](https://img.shields.io/github/license/streimelstefan/tyro)

## 🌟 Overview

`Tyro` is a powerful, interactive terminal user interface (TUI) tool designed for administrators and researchers to efficiently manage, inspect, and fix common issues in DICOM (Digital Imaging and Communications in Medicine) data. Built with Go for high performance, `Tyro` aims to simplify the often complex task of DICOM file manipulation, especially when dealing with large datasets.

## ✨ Key Features

*   **Dual-Pane Interactive TUI:**
    *   **Left Pane:** A hierarchical tree view of the loaded directory, allowing intuitive navigation through DICOM studies and files.
    *   **Right Pane:** A detailed, navigable tree view of DICOM tags (metadata) for the currently selected file, including support for nested sequences.
*   **High Performance Data Handling:**
    *   **Multi-threaded File Discovery:** Leverages Go's goroutines to rapidly scan large directories for DICOM files.
    *   **Lazy DICOM Parsing:** Only essential metadata (DICOM tags) is loaded into memory when a file is selected, ensuring minimal memory footprint and fast responsiveness even with studies containing thousands of images. Pixel data is never loaded unless explicitly requested (a feature for future consideration).
*   **DICOM Issue Resolution (Administrative Focus):** Provides functionalities to identify and rectify common DICOM data inconsistencies and errors.
*   **Intuitive Editing:** Allows users to traverse the DICOM tag tree and modify tag values directly within the TUI.
*   **Cross-Platform Compatibility:** As a Go application, `Tyro` compiles to a single binary, making it easy to deploy across various operating systems.

## 💡 Use Cases

`Tyro` is ideal for situations where you need to:

*   **Fix Corrupt or Non-Compliant DICOM Files:** Address issues like incorrect character sets, missing UIDs, or malformed tags that prevent files from being processed by PACS or viewers.
*   **De-identify / Anonymize DICOM Data:** Quickly strip or pseudonymize patient-identifying information for research or sharing purposes.
*   **Standardize DICOM Headers:** Correct inconsistencies in patient names, IDs, or accession numbers across a study to ensure data integrity.
*   **Deep Inspection of DICOM Metadata:** Explore the full hierarchy of DICOM tags for debugging, validation, or understanding complex private tags.
*   **Batch Operations (Future):** While currently focused on interactive single-file edits, the underlying architecture supports future expansion into applying changes across multiple files or entire studies.
*   **Quality Assurance:** Verify that incoming DICOM data adheres to expected standards before archival.

## 🏗️ Architecture

`Tyro` follows a Model-View-Controller (MVC)-like pattern, optimized for a TUI environment. This separation ensures a clean codebase, high responsiveness, and maintainability.

### 1. Model (Data & State)

The Model encapsulates the application's state and data structures:

*   **`FileTree`**: Represents the directory structure being explored. Each node can be a directory or a DICOM file.
*   **`DicomData`**: Stores the parsed metadata of a single DICOM file. This is loaded lazily, parsing only the DICOM header and tags when a file is selected by the user. Pixel data is not loaded by default.
*   **`TagTree`**: A hierarchical representation of the DICOM tags for the currently selected `DicomData`, mirroring the structure of DICOM sequences.
*   **`AppState`**: A central struct holding the entire application state, including selected file, current focus (file pane or tag pane), and any operational messages.

### 2. View (TUI Rendering)

The View is responsible for rendering the `AppState` to the terminal:

*   **Two Primary Panes**: Utilizes a TUI library (e.g., `tview` or `Bubble Tea`) to render the `FileTree` in the left pane and the `TagTree` of the selected DICOM file in the right pane.
*   **Status Bar**: Displays contextual information, active file path, and operational messages.
*   **Input Fields**: Provides interactive elements for editing tag values.

### 3. Controller (User Input & Logic)

The Controller manages user interactions and orchestrates state changes:

*   **Event Loop**: Continuously listens for keyboard inputs and internal events.
*   **State Updates**: Translates user inputs (e.g., navigation, edit commands) into updates to the `AppState`.
*   **Asynchronous File Discovery**: Spawns concurrent goroutines to walk the specified directory, feeding discovered DICOM file paths back to the main thread via channels, ensuring the UI remains responsive during scans.
*   **DICOM Operations**: Handles parsing DICOM files, modifying tag values, and persisting changes (saving modified DICOM files).

```mermaid
graph TD
    subgraph Controller
        A[User Input (Keyboard)] --> B{Event Loop};
        B --> C[Update AppState Logic];
        D[File Discovery Goroutines] -- Async File Paths via Channel --> C;
        C -- DICOM Operations --> G[DICOM Library];
    end

    subgraph Model
        E[AppState];
        E --> F[FileTree];
        E --> H[Selected DicomData (Lazy Load)];
        H --> I[TagTree (Parsed Tags)];
    end

    subgraph View
        J[TUI Renderer];
        J --> K[Left Pane: FileTree View];
        J --> L[Right Pane: TagTree View];
        J --> M[Status Bar / Input Fields];
    end

    C --> E;
    E -- Renders --> J;
    G -- Updates DicomData --> H;
```

## 🛠️ Technologies & Libraries

*   **GoLang**: The core programming language, chosen for its concurrency primitives, performance, and cross-compilation capabilities.
*   **`go-dicom`**: The primary library for parsing, manipulating, and writing DICOM files.
*   **`tview` / `Bubble Tea`**: (One of these will be chosen as the TUI framework) For building the interactive terminal user interface, handling layouts, and rendering widgets like trees and input fields.
*   **Go's Concurrency Primitives**: `goroutines` and `channels` are leveraged for efficient, non-blocking file discovery and other background tasks.

## 🚀 Getting Started

*(This section will be filled in once the project is ready for building and running.)*

### Prerequisites

*   Go 1.20+

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/tyro.git
cd tyro

# Build the executable
go build -o tyro .
```

### Usage

```bash
# Run Tyro on a directory containing DICOM files
./tyro <path/to/dicom/directory>

# Example:
./tyro ~/dicom_studies
```

*(Further instructions on keybindings and interactive usage will be added here.)*

## 🚧 Common DICOM Compatibility Issues Tyro Aims to Address

Tyro's core mission is to help solve persistent DICOM problems, including:

*   **Incorrect Character Encoding**: Viewing and fixing garbled patient names due to `Specific Character Set` (0008,0005) mismatches (e.g., ISO-IR-100 vs. UTF-8).
*   **Missing or Duplicate UIDs**: Identifying and generating unique `Study Instance UID`, `Series Instance UID`, or `SOP Instance UID` to resolve archiving conflicts.
*   **Inconsistent Patient/Study Attributes**: Standardizing `Patient's Name`, `Patient ID`, `Accession Number`, etc., across an entire study.
*   **Invalid Value Representations (VRs)**: Highlighting and correcting tag values that do not conform to their specified Value Representation.
*   **Private Tag Management**: Inspecting and potentially removing problematic private tags that cause issues with certain viewers or PACS systems.
*   **Anonymization**: Providing a streamlined process to clear or replace sensitive patient information from datasets.

## 🤝 Contributing

We welcome contributions! If you have suggestions, bug reports, or want to contribute code, please open an issue or submit a pull request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.