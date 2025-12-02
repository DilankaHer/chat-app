import { Fragment, useEffect, useState } from "react";
import Login from "./auth/login";
import Room from "./Room";

export interface LoginState {
    isLoginSuccess: boolean;
    userId: string;
}
interface Room {
    roomId: string;
    name: string;
}
function App() {
  const [loginState, setLoginState] = useState<LoginState>({
    isLoginSuccess: false,
    userId: "",
  });
  const [roomId, setRoomId] = useState("");
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [text, setText] = useState("");

  const handleCreateRoom = (e: React.FormEvent) => {
    e.preventDefault();
    const roomName = text;
    fetch("http://localhost:8080/createRoom", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ roomName: roomName }),
    })
      .then(res => {
        if (res.status !== 200) throw new Error("Failed to create room");
        return res.json();
      })
      .then(room => {
        setRooms([...rooms, room]);
      })
      .catch((error) => {
        console.log("Failed to create room", error);
      });
  }

  useEffect(() => {
  fetch("http://localhost:8080/me", {
    credentials: "include"
  })
    .then(res => {
      if (res.status === 200) return res.json();
      throw new Error("Not logged in");
    })
    .then(user => {
      setLoginState({isLoginSuccess: true, userId: user.userId}); // logged in
      setIsLoading(false);
    })
    .catch(() => {
      setLoginState({isLoginSuccess: false, userId: ""}); // not logged in
      setIsLoading(false);
    });
}, []);

useEffect(() => {
  if (loginState.isLoginSuccess) {
    fetch("http://localhost:8080/rooms", {
      credentials: "include"
    })
      .then(res => {
        if (res.status === 200) return res.json();
        throw new Error("Failed to get rooms");
      })
      .then(rooms => {
        setRooms(rooms); // logged in
      })
      .catch(() => setRooms([])); // not logged in
  }
}, [loginState.isLoginSuccess]);

  return (
    <div className="p-40">
      {isLoading ? (<div className="flex justify-center items-center py-10">
        <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
      </div>) :
      (!loginState.isLoginSuccess ? (<Login setLoginState={setLoginState} />) :
      (<Fragment>
        <div className="flex justify-center">
          <span className="text-2xl lg:text-5xl font-mono">Chat App - Room Based WebSocket Demo</span>
        </div>
        <div className="flex flex-row justify-center items-center mt-10">
          <form onSubmit={handleCreateRoom} className="flex flex-row gap-3 mt-4 px-4 py-2 w-full justify-center">
            <input
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Room Name"
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
            <button type="submit" className="text-xl items-center gap-2 rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed">
              Create Room
            </button>
          </form>
        </div>
        <div>
          { roomId !== "" && (
            <Room roomId={roomId} roomName={rooms.find((room) => room.roomId === roomId)?.name || ""} userId={loginState.userId} setRoomId={setRoomId}/>
          )}
          { rooms.length > 0 && roomId === "" && (
            <div className="flex flex-row justify-center items-center p-4 gap-10 mt-20">
            {rooms.map((room) => (
                <button key={room.roomId} className="text-2xl items-center gap-2 rounded-2xl px-20 py-10 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={() => setRoomId(room.roomId)}>{room.name}</button>
            ))}
            </div>
          )}
        </div>
      </Fragment>))}
    </div>
  );
}

export default App;
