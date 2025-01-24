# gowin 

remove server for managing windows server 

### requirements 

- gcc https://sourceforge.net/projects/mingw-w64/files

- gocv https://gocv.io/getting-started/windows/

Step 1: Install Prerequisites
Install Go

Download the latest Windows installer from golang.org.

Run the installer and follow the prompts (ensure "Add to PATH" is checked).

Install OpenCV (Required for GoCV)
GoCV requires OpenCV (v4.6.0 or newer). The easiest way to install OpenCV on Windows is via vcpkg:

bash
Copy
###### Install vcpkg (package manager for C/C++)
git clone https://github.com/microsoft/vcpkg
cd vcpkg
.\bootstrap-vcpkg.bat

###### Install OpenCV with contrib modules
.\vcpkg install opencv[contrib]:x64-windows

###### Integrate vcpkg with Visual Studio (optional but helpful)
.\vcpkg integrate install
Install CMake
Download CMake from cmake.org and add it to your system PATH.

Step 2: Set Environment Variables
Set OPENCV_DIR
Point to the OpenCV installation directory (created by vcpkg):

Copy
OPENCV_DIR = C:\path\to\vcpkg\installed\x64-windows
Add OpenCV to PATH
Add the OpenCV binaries directory to your system PATH:

Copy
PATH = %OPENCV_DIR%\bin;...

### dev 

run `air` command 



