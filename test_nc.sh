#!/bin/bash

echo "=== Testing BBS with nc (netcat) ==="
echo ""

# Start the BBS server in background
echo "Starting BBS server..."
./bbs &
BBS_PID=$!

# Wait for server to start
sleep 2

echo "BBS server started (PID: $BBS_PID)"
echo ""

# Test nc connection
echo "Testing nc connection..."
echo "This will send a few test commands to verify nc works:"
echo ""

# Create a test script with commands
cat > test_commands.txt << 'EOF'
R
testuser
password123
help
rooms
Hello from nc!
quit
EOF

echo "Sending test commands via nc..."
echo ""

# Send commands to nc
timeout 10s nc localhost 3003 < test_commands.txt

echo ""
echo "=== Test completed ==="
echo ""

# Cleanup
kill $BBS_PID 2>/dev/null
rm -f test_commands.txt

echo "BBS server stopped"
echo ""
echo "If you saw the BBS welcome message and help text above,"
echo "then nc is working correctly with the BBS!"
echo ""
echo "To connect manually:"
echo "  nc localhost 3003"
echo ""
echo "To debug connection issues:"
echo "  ./bbs debug"
echo "  nc localhost 3004"