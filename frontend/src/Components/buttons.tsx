import { IconButton, Flex, HStack } from "@chakra-ui/react";
import {
  FaMoon,
  FaSun,
  FaGithub,
} from "react-icons/fa";
import { useColorMode } from "@chakra-ui/react";

export const Buttons = () => {
  const { colorMode, toggleColorMode } = useColorMode();
  return (
    <Flex align="center" justify="end">
      <HStack pt="3" pr="5">
        <IconButton
          colorScheme="purple"
          onClick={() => window.open("https://github.com/cmwaters/blobusign", "_blank")}
          aria-label={`See the repository on GitHub.`}
        >
          <FaGithub />
        </IconButton>
        <IconButton
          onClick={toggleColorMode}
          aria-label={`Switch from ${colorMode} mode`}
        >
          {colorMode === "light" ? <FaSun /> : <FaMoon />}
        </IconButton>
      </HStack>
    </Flex>
  );
};
