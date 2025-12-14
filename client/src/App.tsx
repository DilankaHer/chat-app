import { Fragment, useEffect, useState } from 'react';
import Login from './auth/login';
import Room from './Room';
import { useApi, type ApiRequest } from './utils/apiRequest';

export interface LoginState {
  isLoginSuccess: boolean;
  userId: string;
}
interface Room {
  roomId: string;
  name: string;
}
interface User {
  userId: string;
  username: string;
  email: string;
}
function App() {
  const { apiRequest } = useApi();
  const [loginState, setLoginState] = useState<LoginState>({
    isLoginSuccess: false,
    userId: '',
  });
  const [roomId, setRoomId] = useState('');
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [text, setText] = useState('');

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    const roomName = text;
    const req: ApiRequest = {
      url: '/createRoom',
      method: 'POST',
      body: { roomName: roomName },
    };
    apiRequest<Room>(req)
      .then((response) => {
        console.log(response);
        setRooms([...rooms, response]);
      })
      .catch(() => {
        console.log('failed top create room');
      });
  };

  useEffect(() => {
    console.log('Environment', import.meta.env.MODE);
    if (loginState.isLoginSuccess) return;
    const req: ApiRequest = {
      url: '/me',
      method: 'GET',
      dialogType: 'toast',
    };
    apiRequest<User>(req)
      .then((response) => {
        console.log(response, response.userId);
        setLoginState({ isLoginSuccess: true, userId: response.userId });
        setIsLoading(false);
      })
      .catch(() => {
        setLoginState({ isLoginSuccess: false, userId: '' });
        setIsLoading(false);
        console.log('error');
      });
  }, [loginState.isLoginSuccess]);

  useEffect(() => {
    if (loginState.isLoginSuccess) {
      const req: ApiRequest = {
        url: '/rooms',
        method: 'GET',
      };
      apiRequest<Room[]>(req)
        .then((response) => {
          console.log(response);
          setRooms(response);
        })
        .catch(() => {
          setRooms([]);
        });
    }
  }, [loginState.isLoginSuccess]);

  // const handleError = () => {
  //   console.log("error");
  //   setLoginState({isLoginSuccess: false, userId: ""});
  //   setIsLoading(false);
  // }

  return (
    <div className="flex flex-col gap-3 h-screen overflow-hidden">
      {isLoading ? (
        <div className="flex justify-center items-center py-10">
          <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>
      ) : !loginState.isLoginSuccess ? (
        <Login setLoginState={setLoginState} />
      ) : (
        <Fragment>
          <div className="flex justify-center">
            <span className="text-2xl lg:text-5xl font-mono mt-10">
              Chat App - Room Based WebSocket Demo
            </span>
          </div>
          {roomId === '' && (
            <div className="flex flex-row justify-center items-center">
              <form
                onSubmit={handleCreateRoom}
                className="flex flex-row gap-3 px-4 py-2 w-full justify-center mt-10"
              >
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
                <button
                  type="submit"
                  className="text-xl items-center gap-2 rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Create Room
                </button>
              </form>
            </div>
          )}
          {roomId !== '' && (
            <Room
              roomId={roomId}
              roomName={
                rooms.find((room) => room.roomId === roomId)?.name || ''
              }
              userId={loginState.userId}
              setRoomId={setRoomId}
            />
          )}
          {rooms.length > 0 && roomId === '' && (
            <div className="flex flex-row justify-center items-center p-4 gap-10 mt-20">
              {rooms.map((room) => (
                <button
                  key={room.roomId}
                  className="text-2xl items-center gap-2 rounded-2xl px-20 py-10 font-medium shadow-sm transition bg-linear-to-b from-indigo-600 to-indigo-500 hover:from-indigo-700 hover:to-indigo-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed"
                  onClick={() => setRoomId(room.roomId)}
                >
                  {room.name}
                </button>
              ))}
            </div>
          )}
        </Fragment>
      )}
    </div>
  );
}

export default App;
