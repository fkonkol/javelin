import create from 'zustand';

interface AccountNode {
  id: number;
  username: string;
}

interface SearchState {
  accounts: AccountNode[];
  fetch: (value: string) => Promise<void>;
}

const useStore = create<SearchState>(set => ({
    accounts: {} as AccountNode[],
    fetch: async (value: string) => {
      const response = await fetch(`http://localhost:8000/users?username=${value}`);
      const data = await response.json();
      set({ accounts: data.users });
    }
}));

export default useStore;
