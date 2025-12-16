import React, { Fragment, useEffect, useRef, useState } from 'react';
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
    <Fragment>
      <div className="flex flex-col gap-4 justify-center items-center">
        <div className="flex flex-row gap-3 items-center mt-10">
          <button
            className="text-xl items-center gap-2 rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
            onClick={handleExitRoom}
          >
            Exit Room
          </button>
          <h1 className="lg:text-3xl md:text-2xl text-xl">Connected to Room {roomName}</h1>
        </div>
        <hr className="my-8 h-px bg-neutral-100 dark:bg-neutral-700 border-0 w-full" />
        {/* <h1 className="text-5xl underline"> Messages </h1> */}
      </div>
      <div className="flex flex-col-reverse gap-4 overflow-x-hidden overflow-y-auto lg:text-xl md:text-lg text-base mx-2">
        <ul className="list-none space-y-2 w-full sm:px-25 md:px-50 lg:px-65 xl:px-100 2xl:px-125 mb-5">
          {messages && messages.length > 0 && messages.map((msg, idx) => (
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
      </div>
      <div className='block bg-gray-700 w-full py-4 px-2'>
        <form
          className="flex flex-row gap-3 mb-1"
          onSubmit={handleSubmit}
        >
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
              text-base
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
    </Fragment>
  );
}

export default Room;
