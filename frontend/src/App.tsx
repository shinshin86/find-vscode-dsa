import { useEffect, useState } from 'react';
import './App.css';
import { OpenDir, ProjectInfoList, UpdateWorkspaceRecommendationsIgnore } from "../wailsjs/go/main/App";
import { main } from '../wailsjs/go/models';
import { ProjectInfoList as ProjectInfoListComponent } from './ProjectInfoList';
import { Button, Checkbox, Input, VStack, Box, Text } from '@chakra-ui/react';

function App() {
  const [projectList, setProjectList] = useState<Array<main.ProjectInfo>>([]);
  const [changedProjects, setChangedProjects] = useState<Array<main.ProjectInfo>>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isFiltered, setIsFiltered] = useState(false);
  const [filterText, setFilterText] = useState("");

  const filteredList = projectList.filter((project) => {
    const nameMatch = project.projectName.toLowerCase().includes(filterText.toLowerCase());
    const activeMatch = isFiltered ? project.workspaceRecommendationsIgnore : true;
    return nameMatch && activeMatch;
  });

  const handleCheckboxChange = (changedProject: main.ProjectInfo) => {
    const updateTargetProject = projectList.find(({ vscodeWorkspaceStoragePath }) => vscodeWorkspaceStoragePath === changedProject.vscodeWorkspaceStoragePath);
    if (updateTargetProject) {
      updateTargetProject.workspaceRecommendationsIgnore = changedProject.workspaceRecommendationsIgnore;
    }
    setProjectList(projectList);

    setChangedProjects((prev) => {
      const exists = prev.some((p) => p.vscodeWorkspaceStoragePath === changedProject.vscodeWorkspaceStoragePath);

      if (exists) {
        return prev.filter((p) => p.vscodeWorkspaceStoragePath !== changedProject.vscodeWorkspaceStoragePath);
      } else {
        return [...prev, changedProject];
      }
    });
  };

  const handleSubmit = async () => {
    await UpdateWorkspaceRecommendationsIgnore(changedProjects);
    setChangedProjects([]);
    getProjectInfoList();
  };

  const openDir = (path: string) => OpenDir(path)

  const getProjectInfoList = async () => {
    setIsLoading(true);
    try {
      const response = await ProjectInfoList();
      setProjectList(response);
    } catch (error) {
      console.error(error)
    }
    setIsLoading(false);
  }

  useEffect(() => {
    getProjectInfoList();
  }, [])

  if (isLoading) {
    return <div>loading...</div>
  }

  return (
    <Box id="App" m={4} p={4}>
      <VStack spacing={6}>
        <Text fontSize="3xl" as="b">Ignore Workspace Recommendations List</Text>
        <VStack spacing={4} alignItems="start">
          <Checkbox isChecked={isFiltered} onChange={() => setIsFiltered(!isFiltered)}>Show active items only</Checkbox>
          <Input placeholder="Filter by name" onChange={(e) => setFilterText(e.target.value)} />
        </VStack>
        <Button onClick={() => handleSubmit()} colorScheme="blue" mt={4} isDisabled={changedProjects.length === 0}>
          Submit Changes
        </Button>
        <ProjectInfoListComponent projectInfoList={filteredList} handleCheckboxChange={handleCheckboxChange} openDir={openDir} />
      </VStack>
    </Box>
  )
}

export default App
