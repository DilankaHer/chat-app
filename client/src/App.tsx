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
        setRooms([...rooms, response]);
      })
      .catch(() => {
        console.log('failed top create room');
      });
  };

  useEffect(() => {
    if (loginState.isLoginSuccess) return;
    const req: ApiRequest = {
      url: '/me',
      method: 'GET',
      dialogType: 'toast',
    };
    apiRequest<User>(req)
      .then((response) => {
        setLoginState({ isLoginSuccess: true, userId: response.userId });
        setIsLoading(false);
      })
      .catch(() => {
        setLoginState({ isLoginSuccess: false, userId: '' });
        setIsLoading(false);
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
          setRooms(response);
        })
        .catch(() => {
          setRooms([]);
        });
    }
  }, [loginState.isLoginSuccess]);

  const handleLogout = () => {
    const req: ApiRequest = {
      url: '/logout',
      method: 'POST',
      dialogType: 'toast',
    };
    apiRequest<void>(req)
      .then(() => {
        setLoginState({ isLoginSuccess: false, userId: '' });
        setIsLoading(false);
      })
      .catch(() => {
        setLoginState({ isLoginSuccess: false, userId: '' });
        setIsLoading(false);
      });
  };

  return (
    <div className="flex flex-col gap-3 h-screen overflow-hidden lg:text-xl md:text-lg text-base">
      {isLoading ? (
        <div className="flex justify-center items-center py-10">
          <div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>
      ) : !loginState.isLoginSuccess ? (
        <Login setLoginState={setLoginState} />
      ) : (
        <Fragment>
          <div className='text-end'>
            <button className="items-center rounded-2xl px-4 py-2 m-2 font-medium shadow-sm transition bg-linear-to-b from-red-600 to-red-500 hover:from-red-700 hover:to-red-600 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-red-300 disabled:opacity-50 disabled:cursor-not-allowed"
             onClick={handleLogout}>Logout</button>
          </div>
          <div className="flex justify-center text-center">
            {roomId !== '' ? (
              <span className="lg:text-4xl md:text-3xl text-xl font-bold">
                Connected to {rooms.find((room) => room.roomId === roomId)?.name || ''}
              </span>
            ) : (
              <span className="lg:text-4xl md:text-3xl text-xl font-mono">
                Chat App - Room Based WebSocket Demo
              </span>
            )}
          </div>
          {roomId === '' && (
            <div className="flex flex-row sm:flex-col justify-center items-center">
              <form
                onSubmit={handleCreateRoom}
                className="flex md:flex-row flex-col gap-3 px-4 py-2 w-full justify-center mt-10"
              >
                <input
                  value={text}
                  onChange={(e) => setText(e.target.value)}
                  placeholder="Room Name"
                  className="
                    w-full max-w-md
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
                <div className="text-center">
                  <button
                    type="submit"
                    className="w-fit rounded-2xl px-4 py-2 font-medium shadow-sm transition bg-linear-to-b from-green-600 to-green-500 hover:from-green-800 hover:to-green-800 active:translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-green-300 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Create Room
                  </button>
                </div>
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
            <div className="flex flex-wrap justify-center items-center p-4 gap-4 mt-10">
              {rooms.map((room) => (
                <button
                  key={room.roomId}
                  className="
                    w-fit
                    px-6 py-4 sm:px-10 sm:py-5
                    rounded-2xl
                    font-medium
                    shadow-sm
                    bg-linear-to-b from-indigo-600 to-indigo-500
                  "
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
