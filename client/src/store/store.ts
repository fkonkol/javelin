import create from 'zustand';
import { devtools, persist } from 'zustand/middleware';

interface AuthState {
  loggedIn: boolean;
  login: () => void;
  auth: () => Promise<void>;
}

const useStore = create<AuthState>()(devtools(persist((set) => ({
  loggedIn: false,
  login: () => set({ loggedIn: true }),
  auth: async () => {
    const response = await fetch("http://localhost:8000/auth");
    if (response.status === 200) {
      set({ loggedIn: true });
    } else {
      set({ loggedIn: false });
    }
  }
}))));

export default useStore;
