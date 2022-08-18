/* eslint-disable @typescript-eslint/no-shadow */
import React, {useEffect, useState} from 'react';
import {View, SafeAreaView, Text} from 'react-native';
import WebView from 'react-native-webview';
import StaticServer from '@dr.pogodin/react-native-static-server';
import RNFetchBlob from 'rn-fetch-blob';
var RNFS = require('react-native-fs');

const path = RNFS.MainBundlePath + '/www';
const server = new StaticServer(8080, path);

const App = () => {
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [loaded, setLoaded] = useState<boolean>(false);
  useEffect(() => {
    console.log(RNFS.readDir(path));
    RNFetchBlob.config({path: path + '/image.jpg'})
      .fetch(
        'GET',
        'https://assets.skyweaver.net/aL_BvVlm/webapp/backgrounds/private.jpg',
      )
      .then(res => {
        console.log('file saved to path:', res.path());
        setLoaded(true);
      });
    server.start().then(url => {
      setUrl(url);
    });
  }, []);
  return (
    <SafeAreaView>
      {!loaded && <Text>Downloading and saving image</Text>}
      {url && loaded && (
        <View style={{height: '100%', width: '100%'}}>
          <Text>Serving {url}</Text>
          <WebView style={{flex: 1}} source={{uri: url}} />
        </View>
      )}
    </SafeAreaView>
  );
};

export default App;
