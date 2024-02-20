import axios from 'axios';
import { useState, ChangeEvent } from 'react'; // Import ChangeEvent for typing the event
import { Button, Stack, VStack, Heading, Input, Text } from '@chakra-ui/react'; // Assuming you want to use Text for displaying the file name
import { Buttons } from "./Components/buttons";
import { FaUpload, FaFile, FaSignature } from 'react-icons/fa';

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files ? event.target.files[0] : null;
    if (file) {
      setSelectedFile(file);
      console.log(file); // For demonstration purposes
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
  
    try {
      const response = await axios.post('http://localhost:8080', formData, {
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
        <Heading size="md">Document transparency with Celestia underneath ✨</Heading>
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
        {selectedFile && <Text mt={2} color="gray">Selected file: {selectedFile.name}</Text>}
        <Text>
          Made in space by 👽
        </Text>
      </VStack>
    </>
  );
}

export default App;