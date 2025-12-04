import React, { useEffect, useRef, useState } from "react";
import { useApi, type ApiRequest } from "./utils/apiRequest";

function Room({ roomId, roomName, userId, setRoomId }: { roomId: string, roomName: string, userId: string, setRoomId: (roomId: string) => void }) {
  const { apiRequest } = useApi();
  const [text, setText] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  interface Message {
    roomId: string,
    content: string,
    userId: string,
    username: string,
    isError: boolean,
  }

  useEffect(() => {
    console.log("UserID", userId)
    const ws = new WebSocket(`ws://localhost:8080/joinRoom?roomId=${roomId}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket open for room", roomId);
      const req: ApiRequest = {url: `http://localhost:8080/messages?roomId=${roomId}`, method: "GET"}
      apiRequest(req).then(response => {
        setMessages(response.data);
      })
        .catch(error => {
          console.error("Error fetching messages:", error)
        })
    }
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      setMessages(prevMessages => [
        ...prevMessages,
        { roomId: msg.roomId, content: msg.content, userId: msg.userId, username: msg.username, isError: msg.isError }
      ])
    }
    ws.onclose = () => console.log("WS closed");
    ws.onerror = (err) => console.error("WS error", err);

    return () => ws.close();
  }, [roomId]); // important: re-connect when room changes

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
      setText("");
    }
  };

  const handleExitRoom = () => {
    wsRef.current?.close();
    setRoomId("");
  };

  return (
    <div className="bg-blue-400 dark:bg-blue-950">
      <div className="flex flex-col gap-4 mt-20 justify-center items-center">
        <div className="flex flex-row gap-3 items-center">
        <button className="text-xl items-center gap-2 rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
          type="button" onClick={handleExitRoom}>Exit Room</button>
        <h1 className="text-5xl">Connected to Room {roomName}</h1>
        </div>
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
                flex flex-col
                ${msg.userId === userId ? "text-pink-500 text-right bg-pink-100 ml-auto" : "text-white text-left bg-gray-700 mr-auto"}
              `}
            >
              <span className={`${msg.userId !== userId ? "text-green-500 text-left": "hidden"}`} >{msg.username}</span>
              {msg.content}
            </li>
          ))}
        </ul>
      </div>
  );
}

export default Room;
