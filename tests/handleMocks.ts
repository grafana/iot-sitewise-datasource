import { Page, Response } from '@playwright/test';
import path from 'path';
import { promises as fsPromises, constants as fsConstants } from 'fs';

const mocksFolderPath = path.join(__dirname, 'mocks');

const generateMock = async (mockName: string, mockResponse: Response) => {
  // First, check that the mocks folder exists and if not create it
  try {
    await fsPromises.access(mocksFolderPath, fsConstants.F_OK);
  } catch (err) {
    console.log('Creating mocks folder');
    await fsPromises.mkdir(mocksFolderPath, { recursive: true });
  }

  // then take the response body and write it to a json file in the mocks folder
  try {
    const bufferBody = await mockResponse.body();
    const filePath = `${mocksFolderPath}/${mockName}.json`;
    await fsPromises.writeFile(filePath, bufferBody);
  } catch (err) {
    throw new Error(`Unable to generate mocks for ${mockName}: - ${err}`);
  }
};

// for a given endpoint, subscribes to the response and generates/stores a mock file
const generateMockForEndpoint = async (page: Page, endpoint: string, mockName: string) => {
  page.on('response', async (response) => {
    const url = response.url();
    if (url.includes(endpoint)) {
      generateMock(mockName, response);
    }
  });
};

// for a given endpoint, will replace a response with a mocked response if one exists
const returnMockForEndpoint = async (page: Page, endpoint: string, mockName: string) => {
  await page.route(`**${endpoint}`, async (route) => {
    const mockFiles = await fsPromises.readdir(path.join(__dirname, 'mocks'));
    if (mockFiles.includes(`${mockName}.json`)) {
      const filePath = path.join(__dirname, `mocks/${mockName}.json`);
      const data = await fsPromises.readFile(filePath);
      route.fulfill({ body: data });
    } else {
      throw new Error(`No mock found for ${mockName}`);
    }
  });
};

// will return back either a mocked or live provisioned datasource depending on settings
// also responsible for generating mocks if GENERATE_MOCKS is set to true
// to run tests with mocks: yarn run test:e2e:use-mocks
// to run tests with live endpoints: yarn run test:e2e:use-live-data
// to run tests with live endpoints and generate mocks from the responses: yarn run test:e2e:use-live-data:generate-mocks
// when using live endpoints be sure to create a provisioned datasource in provisioning/datasources/iot-sitewise.e2e.yaml
export const handleMocks = async (page: Page, endpoint: string, testCaseName: string) => {
  if (process.env.GENERATE_MOCKS === 'true') {
    generateMockForEndpoint(page, endpoint, testCaseName);
  }

  if (process.env.USE_MOCKS === 'true') {
    returnMockForEndpoint(page, endpoint, testCaseName);
  }

  let fileName = 'iot-sitewise.e2e.yaml';
  if (process.env.USE_MOCKS === 'true') {
    fileName = 'mock-iot-sitewise.e2e.yaml';
  }

  return {
    fileName: fileName,
    name: testCaseName,
  };
};
