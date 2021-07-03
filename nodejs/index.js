'use strict';

const puppeteer = require('puppeteer');

const result = async (url) => {
    let collectData = {}
    const browser = await puppeteer.launch({headless: true});
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    page.once('load', () => console.log('Page loaded!'));
    page.on('request', (request) => {
        if(['image', 'stylesheet', 'font', 'script'].indexOf(request.resourceType()) !== -1) request.abort();
        else request.continue();
    });
    const response = await page.goto(url, {waitUntil: 'networkidle0'});
    
    // Status
    collectData.status = await response.status()

    // X-Tag-Robot
    let response_headers = response.headers()
    collectData.x_robot_tag = typeof response_headers['x-robots-tag'] !== 'undefined' ? response_headers['x-robots-tag'] : null

    // Meta Robot
    collectData.meta_googlebot = await page.$eval('meta[name=googlebot]', el => el.content).catch(err => null)
    collectData.meta_robot = await page.$eval('meta[name=robots]', el => el.content).catch(err => null)

    await browser.close();

    return collectData
}

let url = 'https://support.google.com/webmasters/answer/1061943'
result(url)
    .then(r => console.log(r))
    .catch(e => console.log('HATA'))
