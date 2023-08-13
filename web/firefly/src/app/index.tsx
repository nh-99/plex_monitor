/**
 *
 * App
 *
 * This component is the skeleton around the actual pages, and should only
 * contain code that should be seen on all pages. (e.g. navigation bar)
 */

import * as React from 'react';
import { Helmet } from 'react-helmet-async';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { GlobalStyle } from 'styles/global-styles';

import { HomePage } from './pages/HomePage/Loadable';
import { Login } from './pages/Login/Loadable';
import { NotFoundPage } from './components/NotFoundPage/Loadable';
import { translations } from 'locales/translations';
import { useTranslation } from 'react-i18next';

export function App() {
  const { t, i18n } = useTranslation();
  const appTitle = t(translations.App.name);
  return (
    <BrowserRouter>
      <Helmet
        titleTemplate={"%s - " + appTitle}
        defaultTitle={appTitle}
        htmlAttributes={{ lang: i18n.language }}
      >
        <meta name="description" content="A Plex (+ Servarr, Ombi) monitoring service" />
      </Helmet>

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
      <GlobalStyle />
    </BrowserRouter>
  );
}