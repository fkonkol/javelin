import React, { ReactNode, useEffect, useState } from 'react';

import { ReactComponent as SearchIcon } from '../assets/img/search.svg';
import useStore from '../store/search';

export const SearchBar: React.FC = () => {
    const [value, setValue] = useState<string>("");
    const [dropdownActive, setDropdownActive] = useState<boolean>(false);

    let accountsFetched = useStore(state => state.accounts);
    const fetchUsers = useStore(state => state.fetch);

    useEffect(() => {
      if (value.length === 0) {
        setDropdownActive(false);
      } else {
        setDropdownActive(true);
        fetchUsers(value);
        console.log(value);
        console.log(accountsFetched);
      }
    }, [value]);

    return (
        <div className="relative bg-white w-full h-12 rounded-md shadow-md grid grid-cols-[3rem_1fr]">
            <div className="flex items-center justify-center">
                <SearchIcon className="w-5 h-5" />
            </div>
            <div>
                <input 
                    className="w-full h-full bg-transparent pr-12 outline-none" 
                    type="text" 
                    placeholder="Search"
                    onChange={(e) => {
                      setValue(e.target.value);
                    }}
                />
            </div>
            {dropdownActive ? (
                <div className="absolute transform translate-y-12 w-full bg-white">
                  {accountsFetched.map((acc) => (
                    <p key={acc.id}>{acc.username}</p>
                  ))}
                </div>
            ) : null}
        </div>
    );
};
