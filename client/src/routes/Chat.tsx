import React from 'react';
import { ReactComponent as Logo } from '../assets/img/avatar.svg';
import { SearchBar } from '../components/SearchBar';

export const Chat: React.FC = () => {
  return (
    <main className="flex w-screen h-screen items-center justify-center bg-slate-100">
      <section className="bg-fuchsia-500 w-96 h-screen py-12 px-8">
        <SearchBar />
      </section>
      <section className="w-[40rem] h-screen">
        <div className="w-full h-full bg-red-600 py-12 px-8">
          <div className="w-full h-full bg-white grid grid-co">
            <div className="w-full h-16 bg-orange-400">
              <Logo className="w-12 h-12" />
            </div>
            <div className="flex items-end">
              <div className="w-full h-16 bg-blue-500 px-8">
                <div className="w-full h-full bg-white">hello</div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
};
