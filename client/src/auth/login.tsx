import React, { useEffect, useState } from "react";
import type { LoginState } from "../App";

function Login({ setLoginState }: { setLoginState: (loginState: LoginState) => void }) {
  interface LoginData {
    email?: string;
    emailUsername?: string;
    password: string;
    username?: string;
  }
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");
  const [isSignup, setIsSignup] = useState(false);

  useEffect(() => {
    console.log(isSignup);
  }, [isSignup]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    let url = ""
    let body: BodyInit | null | undefined = null;
    if (isSignup) {
      const signupData: LoginData = { email, password, username };
      url = "http://localhost:8080/signup"
      body = JSON.stringify(signupData);
    } else {
      const emailUsername = email;
      const loginData: LoginData = { emailUsername, password };
      url = "http://localhost:8080/login"
      body = JSON.stringify(loginData);
    }
    try {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: body,
        credentials: 'include',
      });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      console.log(data);
      setLoginState({
        isLoginSuccess: true,
        userId: data.userId,
      });
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <div className="flex items-center justify-center">
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg max-w-md w-full p-8">
        <h2 className="text-2xl font-bold text-center mb-6">
          {isSignup ? "Register Your Account" : "Login to Your Account"}
        </h2>

        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Email */}
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium"
            >
              {isSignup ? "Email" : "Email or Username"}
            </label>
            <input
              type={isSignup ? "email" : "text"}
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700"
            />
          </div>

          {/* Password */}
          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium "
            >
              Password
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700"
            />
          </div>

          {/* Username */}
          {isSignup && (
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium"
            >
              Username
            </label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-gray-700"
            />
          </div>
          )}

          {/* Submit */}
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
          >
            {isSignup ? "Sign Up" : "Login"}
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm">
            {isSignup? "Already have an account?" : "Don’t have an account?"}
            <a
              className="text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300 font-medium ml-2 cursor-pointer"
              onClick={() => setIsSignup(!isSignup)}
            >
              {isSignup ? "Login" : "Sign Up"}
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
