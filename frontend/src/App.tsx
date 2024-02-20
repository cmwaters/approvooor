import { useState, ChangeEvent } from 'react'; // Import ChangeEvent for typing the event
import { Button, Stack, VStack, Heading, Input, Text } from '@chakra-ui/react'; // Assuming you want to use Text for displaying the file name
import { Buttons } from "./Components/buttons";
import { FaUpload, FaFile } from 'react-icons/fa';

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

  return (
    <>
      <Buttons/>
      <VStack align="start" maxWidth="800px" mx="auto" width="100%" height="100vh">
        <Heading>Blobusign</Heading>
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
          <Button rightIcon={<FaFile />} colorScheme='blue' variant='outline'>
            See Documents
          </Button>
        </Stack>
        {selectedFile && <Text mt={2}>Selected file: {selectedFile.name}</Text>}
        <p className="read-the-docs">
          Made in space by 👽
        </p>
      </VStack>
    </>
  );
}

export default App;