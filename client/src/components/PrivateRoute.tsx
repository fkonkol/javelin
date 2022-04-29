import { AnyArray } from 'immer/dist/internal';
import React, { ReactElement } from 'react';
import { Navigate, Route } from 'react-router-dom';
import { Login } from '../routes/Login';
import useStore from '../store/store';

interface Props {
  Component: any;
  [rest:string]: any;
}

export const PrivateRoute: React.FC<Props> = ({ Component, ...rest }) => {
  const loggedIn = useStore(state => state.loggedIn);
  console.log(loggedIn);

  return loggedIn ? <Component {...rest} /> : <Navigate to="/login" />;
};
