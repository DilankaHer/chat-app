import React, { useEffect, useRef, useState } from "react";

function App() {
  const [text, setText] = useState("");
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Initialize WebSocket only once
    const ws = new WebSocket("ws://localhost:8080/test-message"); // use ws://
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connection is open!");
      ws.send("Hello Server!"); // optional initial message
    };

    ws.onmessage = (event) => {
      console.log("Received from server:", event.data);
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed.");
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    // Clean up on unmount
    return () => {
      ws.close();
    };
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
      setText("");
    } else {
      alert("WebSocket is not connected yet!");
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Type something..."
        />
        <button type="submit">Submit</button>
      </form>
    </div>
  );
}

export default App;
