// i18n.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import translationEN from './locales/en/translation.json';
import translationZH from './locales/zh/translation.json';
import translationJA from './locales/ja/translation.json';
import translationzhHant from './locales/zh-Hant/translation.json';
import translationFR from './locales/fr/translation.json';
import translationKO from './locales/ko/translation.json';
import translationMars from './locales/mars/translation.json';

const resources = {
  en: {
    translation: translationEN,
  },
  zh: {
    translation: translationZH,  
  },
  ja:{
    translation:translationJA,
  },
  zhHant:{
    translation:translationzhHant,
  },
  fr:{
    translation:translationFR,
  },
  ko:{
    translation:translationKO,
  },
  mars:{
    translation:translationMars,
  }
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: 'zh',
    fallbackLng: 'zh',

    keySeparator: false,

    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
