import React, { useEffect, useRef, useState } from "react";

function Room({ id, clientId }: { id: string, clientId: string }) {
  const [text, setText] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  interface Message {
    id: string,
    content: string
    // idx: number
  }

  useEffect(() => {
    console.log("ClientID", clientId)
    const ws = new WebSocket(`ws://localhost:8080/joinRoom?roomId=${id}&clientId=${clientId}`);
    wsRef.current = ws;

    ws.onopen = () => console.log("WebSocket open for room", id);
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      setMessages(prevMessages => [
        ...prevMessages,
        { id: msg.id, content: msg.content }
      ])
    }
    ws.onclose = () => console.log("WS closed");
    ws.onerror = (err) => console.error("WS error", err);

    return () => ws.close();
  }, [id]); // important: re-connect when room changes

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
      setText("");
    }
  };

  return (
    <div className="bg-blue-400 dark:bg-blue-950">
      <div className="flex flex-col gap-4 mt-20 justify-center items-center">
        <h1 className="text-5xl">Connected to Room {id}</h1>
        <form className="flex flex-row gap-3 mt-4 px-4 py-2" onSubmit={handleSubmit}>
          <input
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Say something..."
            className="
              w-full max-w-md
              text-xl
              px-4 py-2
              rounded-xl
              bg-blue-800 dark:bg-white
              placeholder-gray-400 dark:placeholder-gray-500
              shadow-sm
              border border-blue-800 dark:ring-red-900
              focus:outline-none
              focus:ring-2 focus:ring-blue-300 dark:focus:ring-red-900 dark:focus:bg-blue-700
              transition
            "
          />
          <button className="text-xl items-center gap-2 rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"  
          type="submit">Send</button>
        </form>
        <hr className="my-8 h-px bg-neutral-100 dark:bg-neutral-700 border-0 w-full" />
        <h1 className="text-5xl underline"> Messages </h1>
      </div>
      <ul className="list-none mt-10 space-y-2 max-w-5xl mx-auto">
          {messages.map((msg, idx) => (
            <li
              key={idx}
              className={`
                text-2xl 
                px-4 py-2
                rounded-xl 
                w-fit
                max-w-xs
                ${msg.id === clientId ? "text-pink-500 text-right bg-pink-100 ml-auto" : "text-white text-left bg-gray-700 mr-auto"}
              `}
            >
              {msg.content}
            </li>
          ))}
        </ul>
      </div>
  );
}

export default Room;
