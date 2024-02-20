import axios from 'axios';
import { useState, ChangeEvent } from 'react';
import { Button, Stack, VStack, Heading, Input, Text } from '@chakra-ui/react';
import { Buttons } from "./Components/buttons";
import { FaUpload, FaFile, FaSignature } from 'react-icons/fa';

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files ? event.target.files[0] : null;
    if (file) {
      setSelectedFile(file);
      console.log(file);
      console.log("MIME type:", file.type);
    }
  };

  const triggerFileUpload = () => {
    const fileInput = document.getElementById('fileUploadInput') as HTMLInputElement;
    if (fileInput) fileInput.click();
  };

  const sendDocumentToCelestia = async () => {
    if (!selectedFile) {
      alert('No file selected.');
      return;
    }
  
    const formData = new FormData();
    formData.append('document', selectedFile);
    formData.append('mimeType', selectedFile.type);
  
    try {
      const response = await axios.post('http://localhost:8080/submit', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      console.log(response.data);
  
    } catch (error) {
      console.error('Error sending document:', error);
    }
  };

  return (
    <>
      <Buttons/>
      <VStack align="start" maxWidth="800px" mx="auto" width="100%">
        <Heading pb={5}>Blobusign</Heading>
        <Heading size="md" pb={7}>Document transparency with Celestia underneath ✨</Heading>
        <Stack direction='row' spacing={4}>
          <Button leftIcon={<FaUpload />} colorScheme='purple' variant='solid' onClick={triggerFileUpload}>
            Upload Document
          </Button>
          <Input
            type="file"
            id="fileUploadInput"
            style={{ display: 'none' }}
            onChange={handleFileChange}
          />
          {selectedFile && (
            <Button rightIcon={<FaSignature />} colorScheme='green' variant='solid' onClick={sendDocumentToCelestia}>
              Send Document to Celestia
            </Button>
          )}
          <Button rightIcon={<FaFile />} colorScheme='blue' variant='outline'>
            See Documents
          </Button>
        </Stack>
        {selectedFile && <Text mt={2}>Selected file: {selectedFile.name}</Text>}
      </VStack>
    </>
  );
}

export default App;