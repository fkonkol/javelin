import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { ReactComponent as Logo } from '../assets/img/full-gray.svg';

interface UserNode {
  email: string;
  username: string;
  password: string;
}

export const Register: React.FC = () => {
  const [user, setUser] = useState<UserNode>({ email: "", username: "", password: "" });

  let navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const response = await fetch('http://localhost:8000/users/register', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user),
    });

    if (response.status !== 200) {
      return;
    }

    navigate("/chat", { replace: true });
  }

  return (
    <main className="flex items-center justify-center w-screen h-screen bg-slate-100">
      <article className="w-full max-w-md rounded-md shadow-xl p-8">
        <Logo className="w-full h-16 pointer-events-none select-none" />
        <p className="text-2xl font-extralight tracking-wide text-center">Welcome!</p>
        <p className="text-sm text-center text-slate-400 pt-2">
          Already have an account? <Link to="/login" className="text-black">Sign in</Link>
        </p>
        <form onSubmit={handleSubmit}>
          <div className="w-full h-12 rounded-md border-2 border-slate-200 focus-within:border-black my-6">
            <input
              className="bg-transparent w-full h-full px-6 outline-none tracking-wide placeholder-slate-400"
              type="email"
              placeholder="Email"
              onChange={(e) => setUser({ ...user, email: e.target.value })}
            />
          </div>
          <div className="w-full h-12 rounded-md border-2 border-slate-200 focus-within:border-black my-6">
            <input
              className="bg-transparent w-full h-full px-6 outline-none tracking-wide placeholder-slate-400"
              type="text"
              placeholder="Username"
              onChange={(e) => setUser({ ...user, username: e.target.value })}
            />
          </div>
          <div className="w-full h-12 rounded-md border-2 border-slate-200 focus-within:border-black my-6">
            <input
              className="bg-transparent w-full h-full px-6 outline-none tracking-wide placeholder-slate-400"
              type="password"
              placeholder="Password"
              onChange={(e) => setUser({ ...user, password: e.target.value })}
            />
          </div>
          <input
          className="bg-black text-white w-full h-14 rounded-md uppercase font-semibold hover:cursor-pointer hover:bg-gray-800 duration-300 outline-white" 
          type="submit"
          value='Sign up'
          />
        </form>
      </article>
    </main>
  );
}
