import { Box, Checkbox, Text } from '@chakra-ui/react';
import { main } from '../wailsjs/go/models';

export const ProjectInfoList: React.FC<{ projectInfoList: Array<main.ProjectInfo>, handleCheckboxChange: Function, openDir: Function }> = ({ projectInfoList, handleCheckboxChange, openDir }) => (
    <>
        {projectInfoList.map((project, index) => (
            <Box w="full" textAlign={"left"} key={index}>
                <Box m={2}>
                    <Text fontSize="2xl" as="b">{project.projectName}</Text>
                </Box>
                <Box maxWidth="90%" m={2}>
                    <Box>
                        <Text>Project path:</Text>
                    </Box>
                    <Box>
                        <Text as="u" color="blue" style={{ cursor: "pointer" }} onClick={() => openDir(project.projectPath)}>{project.projectPath}</Text>
                    </Box>
                </Box>
                <Box wordBreak="break-word" maxWidth="90%" m={2}>
                    <Box>
                        <Text>VSCode workspace storage path:</Text>
                    </Box>
                    <Box>
                        <Text as="u" color="blue" style={{ cursor: "pointer" }} onClick={() => openDir(project.vscodeWorkspaceStoragePath)}>{project.vscodeWorkspaceStoragePath}</Text>
                    </Box>
                </Box>
                <Checkbox isChecked={project.workspaceRecommendationsIgnore} onChange={() => handleCheckboxChange({ ...project, workspaceRecommendationsIgnore: !project.workspaceRecommendationsIgnore })}>Ignore Workspace Recommendations</Checkbox>
            </Box>
        ))}
    </>
)
