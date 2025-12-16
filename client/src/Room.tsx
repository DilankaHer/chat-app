import React, { useEffect, useRef, useState } from 'react';
import { useApi, type ApiRequest } from './utils/apiRequest';

function Room({
  roomId,
  roomName,
  userId,
  setRoomId,
}: {
  roomId: string;
  roomName: string;
  userId: string;
  setRoomId: (roomId: string) => void;
}) {
  const { apiRequest } = useApi();
  const [text, setText] = useState('');
  const [messages, setMessages] = useState<Message[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  interface Message {
    roomId: string;
    content: string;
    userId: string;
    username: string;
    isError: boolean;
  }

  useEffect(() => {
    if (wsRef.current) return;
    const ws = new WebSocket(`${import.meta.env.VITE_BASE_URL_WS}/joinRoom?roomId=${roomId}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket open for room', roomId);
      const req: ApiRequest = {
        url: `/messages?roomId=${roomId}`,
        method: 'GET',
      };
      apiRequest<Message[]>(req, [])
        .then((response) => {
          setMessages(response);
        })
        .catch((error) => {
          console.error('Error fetching messages:', error);
        });
    };
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, msg]);
    };
    ws.onclose = () => console.log('WS closed');
    ws.onerror = (err) => console.error('WS error', err);

    return () => ws.close();
  }, [roomId]); // important: re-connect when room changes

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
      setText('');
    }
  };

  const handleExitRoom = () => {
    wsRef.current?.close();
    setRoomId('');
  };

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <div className="flex flex-col gap-4 items-center">
        <div className="flex flex-row gap-3 items-center mt-4">
          <button
            className="items-center gap-2 rounded-2xl md:px-4 md:py-2 px-2 py-1 font-medium shadow-sm transition bg-linear-to-b from-red-600 to-red-500 hover:from-red-700 hover:to-red-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-red-300 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
            onClick={handleExitRoom}
          >
            Exit Room
          </button>
        </div>
      <hr className="my-8 h-px bg-neutral-100 dark:bg-neutral-700 border-0 w-full" />
      </div>
      <div className="flex-1 overflow-y-auto overflow-x-hidden px-2">
        <ul className="list-none space-y-2 w-full sm:px-25 md:px-50 lg:px-65 xl:px-100 2xl:px-125 mb-5">
          {messages.map((msg, idx) => (
            <li
              key={idx}
              className={`
                px-4 py-2
                rounded-xl
                w-fit
                max-w-xs
                flex flex-col
                whitespace-pre-wrap
                break-words
                overflow-wrap-anywhere
                ${
                  msg.userId === userId
                    ? 'text-pink-500 text-right bg-pink-100 ml-auto'
                    : 'text-white text-left bg-gray-700 mr-auto'
                }
              `}
            >
              <span
                className={`${
                  msg.userId !== userId ? 'text-green-500 text-left' : 'hidden'
                }`}
              >
                {msg.username}
              </span>
              {msg.content}
            </li>
          ))}
        </ul>
        <div ref={messagesEndRef} />
      </div>
      <div className="bg-gray-700 w-full px-2 py-3">
        <form className="flex items-center gap-3" onSubmit={handleSubmit}>
          <textarea
            rows={1}
            value={text}
            onChange={(e) => {
              setText(e.target.value);
              // e.target.style.height = "auto";
              // e.target.style.height = e.target.scrollHeight + "px";
            }}
            placeholder="Say something..."
            className="
              w-full
              px-4 py-2
              leading-normal
              resize-none
              outline-none
            "
          />
          <button
            type="submit"
            className="
              px-4 py-2
              leading-normal
              rounded-2xl
              font-medium
              shadow-sm
              bg-indigo-600 hover:bg-indigo-700
              mr-4
            "
          >
            Send
          </button>
        </form>
      </div>
    </div>
  );
}

export default Room;
