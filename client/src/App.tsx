import { useState } from "react";
import Room from "./Room";
function App() {
  const [roomId, setRoomId] = useState<string | null>(null);
  const clientId = crypto.randomUUID();

  return (
    <div className="bg-blue-400 dark:bg-blue-950 p-40 min-h-screen text-black dark:text-white">
      <div className="flex justify-center">
        <span className=" text-7xl font-mono">Chat App - Room Based WebSocket Demo</span>
      </div>
      <div>
        {!roomId ? (
          <div className="flex flex-row justify-center items-center p-4 gap-10 mt-90">
            <button className="text-6xl items-center gap-2 rounded-2xl px-20 py-10 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
          onClick={() => setRoomId("1")}>Room 1</button>
            <button className="text-6xl items-center gap-2 rounded-2xl px-20 py-10 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed" 
            onClick={() => setRoomId("2")}>Room 2</button>
          </div>
        ) : (
          <Room id={roomId} clientId={clientId} />
        )}
      </div>
    </div>
  );
}

export default App;
