import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { Login } from './routes/Login';
import { Register } from './routes/Register';
import { PrivateRoute } from './components/PrivateRoute';
import { Chat } from './routes/Chat';
import Home from './routes/Home';
import { PageNotFound } from './routes/PageNotFound';

const App: React.FC = () => {
  return (
    <BrowserRouter>
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="login" element={<Login />} />
      <Route path="register" element={<Register />} />
      <Route path="chat" element={
        <PrivateRoute>
          <Chat />
        </PrivateRoute>
      } />
      <Route path="*" element={<PageNotFound />} />
    </Routes>
  </BrowserRouter> 
  );
}

export default App;
