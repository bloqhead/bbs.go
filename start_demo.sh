#!/bin/bash

echo "=== Enhanced BBS Demo Script ==="
echo ""

# Build the BBS server
echo "Building BBS server..."
go build -o bbs
if [ $? -ne 0 ]; then
    echo "Failed to build BBS server"
    exit 1
fi

# Build the admin tool
echo "Building admin tool..."
cd cmd/admin && go build -o ../../admin . && cd ../..
if [ $? -ne 0 ]; then
    echo "Failed to build admin tool"
    exit 1
fi

echo "✓ Build complete!"
echo ""

echo "=== Quick Start Guide ==="
echo ""
echo "1. Start the BBS server:"
echo "   ./bbs"
echo ""
echo "2. In another terminal, connect with telnet:"
echo "   telnet localhost 3003"
echo ""
echo "3. To manage the BBS (MOTD, rooms, users):"
echo "   ./admin"
echo ""
echo "=== Multiple Users Demo ==="
echo "You can open multiple terminal windows and connect multiple users:"
echo "   Terminal 1: telnet localhost 3003"
echo "   Terminal 2: telnet localhost 3003"
echo "   Terminal 3: telnet localhost 3003"
echo ""
echo "Each user can register, join rooms, and chat in real-time!"
echo ""

# Check if tmux is available for advanced demo
if command -v tmux &> /dev/null; then
    echo "=== Advanced Demo (with tmux) ==="
    echo "Run this to start a tmux session with server and multiple clients:"
    echo "   ./demo_tmux.sh"
    echo ""
fi

echo "Press any key to start the BBS server..."
read -n 1 -s

echo "Starting BBS server on port 3003..."
echo "Connect with: telnet localhost 3003"
echo "Press Ctrl+C to stop the server"
echo ""

./bbs