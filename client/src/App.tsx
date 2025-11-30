import { Fragment, useEffect, useState } from "react";
import Login from "./auth/login";
import Room from "./Room";

export interface LoginState {
    isLoginSuccess: boolean;
    userId: string;
}
interface Room {
    id: string;
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

  useEffect(() => {
  fetch("http://localhost:8080/me", {
    credentials: "include"
  })
    .then(res => {
      if (res.status === 200) return res.json();
      throw new Error("Not logged in");
    })
    .then(user => {
      setLoginState({isLoginSuccess: true, userId: user.id}); // logged in
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
          <span className=" text-7xl font-mono">Chat App - Room Based WebSocket Demo</span>
        </div>
        <div>
          { roomId !== "" && (
            <Room roomId={roomId} roomName={rooms.find((room) => room.id === roomId)?.name || ""} userId={loginState.userId} setRoomId={setRoomId}/>
          )}
          { rooms.length > 0 && roomId === "" && (
            <div className="flex flex-row justify-center items-center p-4 gap-10 mt-90">
            {rooms.map((room) => (
                <button key={room.id} className="text-6xl items-center gap-2 rounded-2xl px-20 py-10 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={() => setRoomId(room.id)}>Room {room.name}</button>
            ))}
            </div>
          )}
        </div>
      </Fragment>))}
    </div>
  );
}

export default App;
