/* eslint-disable @typescript-eslint/no-shadow */
import React, {useEffect, useState} from 'react';
import {View, SafeAreaView, Text, NativeModules} from 'react-native';
import WebView from 'react-native-webview';
import StaticServer from '@dr.pogodin/react-native-static-server';
import RNFetchBlob from 'rn-fetch-blob';
var RNFS = require('react-native-fs');
const AppWebServer = NativeModules.AppWebServer;

const path = RNFS.MainBundlePath + '/www';
const documentPath = RNFS.DocumentDirectoryPath + '/www';
const imagePath = documentPath + '/image.jpg';

let server;

const App = () => {
  const startServer = () => {
    server = new StaticServer(8080, documentPath);
    server.start().then(url => {
      setUrl(url);
    });
  };

  const startGoServer = () => {
    AppWebServer.start(path).then((url: string) => {
      setUrl(url);
    });
  };

  const [isUsingCached, setIsUsingCached] = useState(false);
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [loaded, setLoaded] = useState<boolean>(false);
  useEffect(() => {
    RNFS.exists(documentPath).then((exists: boolean) => {
      if (!exists) {
        console.log(path + ' to ' + documentPath);
        RNFS.copyFile(path, documentPath);
      } else {
        console.log('documentPath already exists, using it');
      }
    });

    RNFS.exists(imagePath).then((exists: boolean) => {
      if (exists) {
        console.log('image already exists, using it');
        setLoaded(true);
        setIsUsingCached(true);
        // startServer();
        startGoServer();
      } else {
        console.log('fetching image');
        RNFetchBlob.config({path: imagePath})
          .fetch(
            'GET',
            'https://assets.skyweaver.net/aL_BvVlm/webapp/backgrounds/private.jpg',
          )
          .then(res => {
            console.log('file saved to path:', res.path());
            setLoaded(true);
            // startServer();
            startGoServer();
          });
      }
    });
  }, []);
  return (
    <SafeAreaView>
      {!loaded && <Text>Downloading and saving image</Text>}
      {url && loaded && (
        <View style={{height: '100%', width: '100%', backgroundColor: 'white'}}>
          <Text>Serving {url}</Text>
          <Text>Path: {documentPath}</Text>
          <Text>Using cached: {isUsingCached.toString()}</Text>
          <WebView
            style={{flex: 1}}
            source={{uri: url}}
            allowFileAccess
            allowingReadAccessToURL={url}
          />
        </View>
      )}
    </SafeAreaView>
  );
};

export default App;
