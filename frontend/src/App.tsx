import axios from 'axios';
import { useState, ChangeEvent } from 'react';
import { Button, Stack, VStack, Heading, Input, Text, useToast } from '@chakra-ui/react';
import { Buttons } from "./Components/buttons";
import { FaUpload, FaSignature, FaDownload } from 'react-icons/fa';

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const toast = useToast();
  const [fileID, setFileID] = useState<string | null>(null);
  const [inputFileID, setInputFileID] = useState<string>('');
  const [fileBytes, setFileBytes] = useState<string | null>(null);

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files ? event.target.files[0] : null;
    if (file) {
      setSelectedFile(file);
    }
  };

  const triggerFileUpload = () => {
    const fileInput = document.getElementById('fileUploadInput') as HTMLInputElement;
    if (fileInput) fileInput.click();
  };

  const handleInputFileIDChange = (event: ChangeEvent<HTMLInputElement>) => {
    setInputFileID(event.target.value);
  };

  const getFileByID = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/get?id=${inputFileID}`);
      setFileBytes(response.data);
      toast({
        title: 'File retrieved successfully.',
        description: 'File bytes have been fetched.',
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
    } catch (error) {
      console.error('Error fetching file:', error);
      toast({
        title: 'Error fetching file.',
        description: "There was an error fetching the file. Please check the file ID.",
        status: 'error',
        duration: 9000,
        isClosable: true,
      });
    }
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
      setFileID(response.data);
      toast({
        title: 'Document sent successfully.',
        description: `Response: ${response.data}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
  
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error('Axios error:', error.message);
        if (error.response) {
          console.error('Response data:', error.response.data);
          console.error('Response status:', error.response.status);
        }
      } else {
        console.error('Unexpected error:', error);
      }
      toast({
        title: 'Error sending document.',
        description: "There was an error uploading the document. Please start your server.",
        status: 'error',
        duration: 9000,
        isClosable: true,
      });
    }
  };

  return (
    <>
      <Buttons/>
      <VStack align="start" maxWidth="800px" mx={{ base: "6", sm: "8", md: "10", lg: "auto" }} width="100%">
        <Heading pb={5}>BlobuSign</Heading>
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
          {/* <Button rightIcon={<FaFile />} colorScheme='blue' variant='outline'>
            Get Documents
          </Button> */}
        </Stack>
        <VStack>
          {selectedFile && <Text mt={2}>Selected file: {selectedFile.name}</Text>}
          {fileID && <Text mt={2}>Uploaded File ID: {fileID}</Text>}
          { fileID && <Input
          placeholder="Enter File ID"
          value={inputFileID}
          onChange={handleInputFileIDChange}
          />}
          { fileID && <Button leftIcon={<FaDownload />} colorScheme='orange' variant='solid' onClick={getFileByID}>
            Get File
          </Button>}
          {fileBytes && <Text mt={2}>File Bytes: {fileBytes}</Text>}
        </VStack>
      </VStack>
    </>
  );
}

export default App;