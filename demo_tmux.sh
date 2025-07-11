#!/bin/bash

# Enhanced BBS tmux demo script
# This creates a tmux session with the BBS server and multiple client connections

SESSION_NAME="bbs-demo"

echo "=== Enhanced BBS Tmux Demo ==="
echo ""

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo "tmux is not installed. Please install tmux to use this demo."
    echo "On Ubuntu/Debian: sudo apt install tmux"
    echo "On macOS: brew install tmux"
    exit 1
fi

# Build everything first
echo "Building BBS components..."
go build -o bbs
cd cmd/admin && go build -o ../../admin . && cd ../..

# Kill existing session if it exists
tmux kill-session -t $SESSION_NAME 2>/dev/null

echo "Starting tmux demo session..."
echo "This will create a session with:"
echo "  - BBS Server (top pane)"
echo "  - Three client connections"
echo "  - Admin tool (bottom right)"
echo ""

# Create new session and split it
tmux new-session -d -s $SESSION_NAME

# Rename the first window
tmux rename-window -t $SESSION_NAME:0 'BBS-Demo'

# Split into multiple panes
tmux split-window -t $SESSION_NAME:0 -h  # Split horizontally (left/right)
tmux split-window -t $SESSION_NAME:0.0 -v  # Split left pane vertically (server/client1)
tmux split-window -t $SESSION_NAME:0.2 -v  # Split right pane vertically (client2/admin)

# Now we have 4 panes:
# 0: Server (top-left)
# 1: Client 1 (bottom-left) 
# 2: Client 2 (top-right)
# 3: Admin tool (bottom-right)

# Start the BBS server in pane 0
tmux send-keys -t $SESSION_NAME:0.0 './bbs' Enter

# Wait a moment for server to start
sleep 2

# Setup client connections (they'll wait for user input)
tmux send-keys -t $SESSION_NAME:0.1 'echo "=== Client 1 ===" && echo "Register with: R, then choose username/password" && telnet localhost 3003' Enter
tmux send-keys -t $SESSION_NAME:0.2 'echo "=== Client 2 ===" && echo "Register with: R, then choose username/password" && telnet localhost 3003' Enter

# Setup admin tool
tmux send-keys -t $SESSION_NAME:0.3 'echo "=== Admin Tool ===" && echo "Commands: motd, room list, users, help" && ./admin'

# Resize panes for better visibility
tmux resize-pane -t $SESSION_NAME:0.0 -y 12  # Server pane height
tmux resize-pane -t $SESSION_NAME:0.1 -y 12  # Client 1 pane height

# Add some helpful text overlays
tmux send-keys -t $SESSION_NAME:0.0 '' # Just to position cursor

# Attach to the session
echo "Attaching to tmux session..."
echo ""
echo "=== Tmux Controls ==="
echo "  Ctrl+B then:"
echo "    Arrow keys - Switch between panes"
echo "    D - Detach from session (server keeps running)"
echo "    X - Kill current pane"
echo "    & - Kill entire session"
echo ""
echo "=== Demo Instructions ==="
echo "1. Register users in Client 1 and Client 2 panes"
echo "2. Try commands like: help, rooms, join Tech, msg Hello!"
echo "3. Use the Admin tool to manage MOTD and rooms"
echo "4. Watch real-time chat between clients!"
echo ""

tmux attach-session -t $SESSION_NAME