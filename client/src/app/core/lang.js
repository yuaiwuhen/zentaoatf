import {app} from "electron";
import path from "path";
const fs = require('fs');

import {DEBUG} from '../utils/consts';
import {logInfo} from "../utils/log";
import LangHelper from '../utils/lang';

const langHelper = new LangHelper();

export const initLang = () => {
    let langName = app.getLocale()

    langName = langName.toLowerCase()
    if (langName.startsWith('zh-')) {
        langName = 'zh-cn';
    } else {
        langName = 'en';
    }

    logInfo(`langName=${langName}`)
    loadLanguage(langName)
};

export const loadLanguage = (langName) => {
    if (!langName) {
        return
    }

    if (langName !== langHelper.name) {
        const langData = loadLangData(langName)
        langHelper.change(langName, langData);
    }
};

const loadLangData = (langName) => {
    let pth = `lang/${langName}.json`
    if (!DEBUG) {
        pth = path.join(process.resourcesPath, pth)
    }

    logInfo(`load language res from ${pth}`)

    const buf = fs.readFileSync(pth)
    const data = JSON.parse(buf.toString())
    return data
}

if (DEBUG) {
    global.$lang = langHelper;
}

export default langHelper;
