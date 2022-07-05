import React from 'react';
import { ListPage } from './pages/ListPage/ListPage';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { UploadPage } from './pages/UploadPage/UploadPage';
import { SigninPage } from './pages/SigninPage/SigninPage';
import { AuthContextProvider } from './auth/AuthContextProvider';
import { ContentPage } from './pages/ContentPage/ContentPage';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <AuthContextProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<ListPage />} />
            <Route path="/upload" element={<UploadPage />} />
            <Route path="/signin" element={<SigninPage />} />
            <Route path="/file/:name" element={<ContentPage />} />
          </Routes>
        </BrowserRouter>
      </AuthContextProvider>
    </ThemeProvider>
  );
}

export default App;
