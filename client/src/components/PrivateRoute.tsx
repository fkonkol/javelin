import React, { useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import useStore from '../store/store';

interface Props {
  children: JSX.Element;
}

export const PrivateRoute: React.FC<Props> = ({ children }) => {
  const auth = useStore(state => state.auth);
  const loggedIn = useStore(state => state.loggedIn);

  useEffect(() => {
    const checkAuth = async () => await auth();
    checkAuth();
  }, [auth]);

  if (!loggedIn) {
    return <Navigate to="/login" />;
  }
  return children;
};
