package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	backendmodule "github.com/sohel902833/go_super_cli/src/backend-module"
	"github.com/sohel902833/go_super_cli/src/types"
	"github.com/spf13/cobra"
)


var (
	configFile string
	dryRun     bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "super",
	Short: "Super CLI - Dynamic CRUD Code Generator",
	Long:  `A powerful CLI tool that generates CRUD boilerplate code dynamically based on configurable instructions.`,
}

var createCmd = &cobra.Command{
	Use:   "create [bm|fm]",
	Short: "Create a new module",
	Long:  `Create a new backend module (bm) or frontend module (fm) interactively.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleType := args[0]
		if moduleType != "bm" && moduleType != "fm" {
			fmt.Println("Error: module type must be 'bm' or 'fm'")
			return
		}
		handleCreate(moduleType)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [filepath]",
	Short: "Bulk create modules from JSON file",
	Long:  `Upload a JSON file containing module definitions and create multiple modules at once.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		handleBulkUpload(filepath)
	},
}

var initCmd = &cobra.Command{
	Use:   "init [bp|fp]",
	Short: "Initialize a new project with base structure",
	Long:  `Create a new backend project (bp) or frontend project (fp) interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectType := args[0]
		if projectType != "bp" && projectType != "fp" {
			fmt.Println("Error: Project type must be 'bp' or 'fp'")
			return
		}
		handleInit(projectType)
	},
}

// var configCmd = &cobra.Command{
// 	Use:   "config",
// 	Short: "Manage configuration templates",
// 	Long:  `Load, save, and manage custom template configurations.`,
// }

// var loadConfigCmd = &cobra.Command{
// 	Use:   "load [filepath]",
// 	Short: "Load custom configuration from file",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		handleLoadConfig(args[0])
// 	},
// }

// var exportConfigCmd = &cobra.Command{
// 	Use:   "export [filepath]",
// 	Short: "Export current configuration to file",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		handleExportConfig(args[0])
// 	},
// }

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(initCmd)
	// rootCmd.AddCommand(configCmd)

	// configCmd.AddCommand(loadConfigCmd)
	// configCmd.AddCommand(exportConfigCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Custom config file path")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview changes without creating files")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func handleCreate(moduleType string) {
	var moduleName string
	var fields string

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter module name: ")
	moduleName, _ = reader.ReadString('\n')
	moduleName = strings.TrimSpace(moduleName)

	if moduleName == "" {
		fmt.Println("Error: module name is required")
		return
	}

	fmt.Print("Enter fields (optional, format: name@S@R,email@S@R): ")
	fields, _ = reader.ReadString('\n')
	fields = strings.TrimSpace(fields)

	module := types.Module{
		ModuleName:      moduleName,
		ModelProperties: fields,
	}

	if dryRun {
		fmt.Println("\n🔍 DRY RUN MODE - No files will be created\n")
	}

	generateModule(moduleType, module)
	
	if !dryRun {
		fmt.Printf("\n✓ Module '%s' created successfully!\n", moduleName)
	}
}

func handleBulkUpload(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var modules []types.Module
	if err := json.Unmarshal(data, &modules); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	if dryRun {
		fmt.Println("\n🔍 DRY RUN MODE - No files will be created\n")
	}

	fmt.Printf("Creating %d modules...\n\n", len(modules))
	for i, module := range modules {
		fmt.Printf("[%d/%d] Creating module: %s\n", i+1, len(modules), module.ModuleName)
		generateModule("bm", module)
		fmt.Println()
	}
	
	if !dryRun {
		fmt.Println("✓ Bulk creation completed!")
	}
}

func handleInit(projectType string) {
	fmt.Println("🚀 Initializing new project...")

	if(projectType=="fp"){
		 fmt.Println("Frontend project initialization is not supported yet!");
		 return;
	}

	instructions,updates:=backendmodule.GetInitProjectInstructions();
	replacements:=map[string]string{
			"PROJECT_NAME":"faltu",
		}
	fmt.Println("📁 Creating files:")
	for _, instruction := range instructions {
		
		filePath := applyReplacements(instruction.FilePath, replacements)
		fileContent := applyReplacements(instruction.Content, replacements)

		if verbose && instruction.Description != "" {
			fmt.Printf("  ℹ %s\n", instruction.Description)
		}

		if dryRun {
			fmt.Printf("  [DRY RUN] Would create: %s\n", filePath)
		} else {
			if err := createFile(filePath, fileContent); err != nil {
				fmt.Printf("  ✗ Error creating %s: %v\n", filePath, err)
				continue
			}
			fmt.Printf("  ✓ Created: %s\n", filePath)
		}
	}

	// Update existing files
	if len(updates) > 0 {
		fmt.Println("\n🔧 Updating files:")
		for _, update := range updates {
			filePath := applyReplacements(update.FilePath, replacements)
			placeholder := applyReplacements(update.Placeholder, replacements)
			content := applyReplacements(update.Content, replacements)

			if verbose && update.Description != "" {
				fmt.Printf("  ℹ %s\n", update.Description)
			}

			if dryRun {
				fmt.Printf("  [DRY RUN] Would update: %s (at placeholder: %s)\n", filePath, placeholder)
			} else {
				if err := updateFile(filePath, placeholder, content, update.Position, update.CreateIfNotExists); err != nil {
					fmt.Printf("  ✗ Error updating %s: %v\n", filePath, err)
					continue
				}
				fmt.Printf("  ✓ Updated: %s\n", filePath)
			}
		}
	}


// 	baseStructure := []string{
// 		"src/modules",
// 		"src/features",
// 		"src/config",
// 		"src/utils",
// 		"src/types",
// 		"src/middlewares",
// 	}

// 	for _, dir := range baseStructure {
// 		if dryRun {
// 			fmt.Printf("  [DRY RUN] Would create: %s/\n", dir)
// 		} else {
// 			if err := os.MkdirAll(dir, 0755); err != nil {
// 				fmt.Printf("  ✗ Error creating %s: %v\n", dir, err)
// 				continue
// 			}
// 			fmt.Printf("  ✓ Created: %s/\n", dir)
// 		}
// 	}

// 	// Create sample config file
// 	sampleConfig := getDefaultConfig()
// 	configJSON, _ := json.MarshalIndent(sampleConfig, "", "  ")
	
// 	if dryRun {
// 		fmt.Println("\n  [DRY RUN] Would create: super.config.json")
// 	} else {
// 		if err := os.WriteFile("super.config.json", configJSON, 0644); err != nil {
// 			fmt.Printf("  ✗ Error creating config: %v\n", err)
// 		} else {
// 			fmt.Println("\n  ✓ Created: super.config.json")
// 		}
// 	}

// 	// Create main route file with placeholders
// 	mainRouteContent := `import { Router } from 'express';

// const router = Router();

// // SUPER_CLI_ROUTE_IMPORTS
// // Auto-generated route imports will be added here

// // SUPER_CLI_ROUTE_DEFINITIONS
// // Auto-generated route definitions will be added here

// export default router;
// `

// 	if dryRun {
// 		fmt.Println("  [DRY RUN] Would create: src/routes/index.ts")
// 	} else {
// 		if err := createFile("src/routes/index.ts", mainRouteContent); err != nil {
// 			fmt.Printf("  ✗ Error creating routes file: %v\n", err)
// 		} else {
// 			fmt.Println("  ✓ Created: src/routes/index.ts")
// 		}
// 	}

	fmt.Println("\n✨ Project initialized! You can now use 'super create' to generate modules.")
}

// func handleLoadConfig(filepath string) {
// 	data, err := os.ReadFile(filepath)
// 	if err != nil {
// 		fmt.Printf("Error reading config file: %v\n", err)
// 		return
// 	}

// 	var config types.ProjectConfig
// 	if err := json.Unmarshal(data, &config); err != nil {
// 		fmt.Printf("Error parsing config: %v\n", err)
// 		return
// 	}

// 	currentConfig = &config
// 	fmt.Printf("✓ Loaded configuration: %s (v%s)\n", config.Name, config.Version)
// 	fmt.Printf("  - %d file instructions\n", len(config.FileInstructions))
// 	fmt.Printf("  - %d update instructions\n", len(config.UpdateInstructions))
// }

// func handleExportConfig(filepath string) {
// 	config := getDefaultConfig()
// 	configJSON, err := json.MarshalIndent(config, "", "  ")
// 	if err != nil {
// 		fmt.Printf("Error creating config: %v\n", err)
// 		return
// 	}

// 	if err := os.WriteFile(filepath, configJSON, 0644); err != nil {
// 		fmt.Printf("Error writing config file: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✓ Configuration exported to: %s\n", filepath)
// }

func generateModule(moduleType string, module types.Module) {
	fields := parseFields(module.ModelProperties)
	replacements := buildReplacements(module.ModuleName, fields)

	// // Load custom config if specified
	// var config *ProjectConfig
	// if configFile != "" {
	// 	handleLoadConfig(configFile)
	// 	config = currentConfig
	// }

	// Get instructions
	// var instructions []types.FileInstruction
	// var updates []types.UpdateInstruction

	// if config != nil {
	// 	instructions = config.FileInstructions
	// 	updates = config.UpdateInstructions
	// } else {
		
	// }
	instructions,updates := getInstructions(moduleType)

	// Create files
	fmt.Println("📁 Creating files:")
	for _, instruction := range instructions {
		filePath := applyReplacements(instruction.FilePath, replacements)
		fileContent := applyReplacements(instruction.Content, replacements)

		if verbose && instruction.Description != "" {
			fmt.Printf("  ℹ %s\n", instruction.Description)
		}

		if dryRun {
			fmt.Printf("  [DRY RUN] Would create: %s\n", filePath)
		} else {
			if err := createFile(filePath, fileContent); err != nil {
				fmt.Printf("  ✗ Error creating %s: %v\n", filePath, err)
				continue
			}
			fmt.Printf("  ✓ Created: %s\n", filePath)
		}
	}

	// Update existing files
	if len(updates) > 0 {
		fmt.Println("\n🔧 Updating files:")
		for _, update := range updates {
			filePath := applyReplacements(update.FilePath, replacements)
			placeholder := applyReplacements(update.Placeholder, replacements)
			content := applyReplacements(update.Content, replacements)

			if verbose && update.Description != "" {
				fmt.Printf("  ℹ %s\n", update.Description)
			}

			if dryRun {
				fmt.Printf("  [DRY RUN] Would update: %s (at placeholder: %s)\n", filePath, placeholder)
			} else {
				if err := updateFile(filePath, placeholder, content, update.Position, update.CreateIfNotExists); err != nil {
					fmt.Printf("  ✗ Error updating %s: %v\n", filePath, err)
					continue
				}
				fmt.Printf("  ✓ Updated: %s\n", filePath)
			}
		}
	}
}

func parseFields(fieldsStr string) []types.Field {
	if fieldsStr == "" {
		return []types.Field{}
	}

	fieldParts := strings.Split(fieldsStr, ",")
	fields := make([]types.Field, 0, len(fieldParts))

	for _, part := range fieldParts {
		segments := strings.Split(strings.TrimSpace(part), "@")
		if len(segments) >= 2 {
			field := types.Field{
				Name:     segments[0],
				Type:     segments[1],
				Required: len(segments) >= 3 && segments[2] == "R",
			}
			fields = append(fields, field)
		}
	}

	return fields
}

func getInstructions(moduleType string) ([]types.FileInstruction,[]types.UpdateInstruction) {
	if moduleType == "bm" {
		return backendmodule.GetCreateBackendModuleInstructions()
	}
	return []types.FileInstruction{},[]types.UpdateInstruction{}
	// return getFrontendInstructions()
}

// func getBackendInstructions() []FileInstruction {
// 	return []FileInstruction{
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.controller.ts",
// 			Description: "Creating controller with CRUD operations",
// 			Content: `import { Request, Response } from 'express';
// import { {{PASCAL_CASE_MODULE_NAME}}Service } from './{{LOWER_CASE_MODULE_NAME}}.service';

// export class {{PASCAL_CASE_MODULE_NAME}}Controller {
//   private service: {{PASCAL_CASE_MODULE_NAME}}Service;

//   constructor() {
//     this.service = new {{PASCAL_CASE_MODULE_NAME}}Service();
//   }

//   async create(req: Request, res: Response) {
//     try {
//       const data = await this.service.create(req.body);
//       res.status(201).json(data);
//     } catch (error) {
//       res.status(500).json({ error: error.message });
//     }
//   }

//   async findAll(req: Request, res: Response) {
//     try {
//       const data = await this.service.findAll();
//       res.status(200).json(data);
//     } catch (error) {
//       res.status(500).json({ error: error.message });
//     }
//   }

//   async findOne(req: Request, res: Response) {
//     try {
//       const data = await this.service.findOne(req.params.id);
//       res.status(200).json(data);
//     } catch (error) {
//       res.status(404).json({ error: error.message });
//     }
//   }

//   async update(req: Request, res: Response) {
//     try {
//       const data = await this.service.update(req.params.id, req.body);
//       res.status(200).json(data);
//     } catch (error) {
//       res.status(500).json({ error: error.message });
//     }
//   }

//   async delete(req: Request, res: Response) {
//     try {
//       await this.service.delete(req.params.id);
//       res.status(204).send();
//     } catch (error) {
//       res.status(500).json({ error: error.message });
//     }
//   }
// }
// `,
// 		},
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.service.ts",
// 			Description: "Creating service layer with business logic",
// 			Content: `import { {{PASCAL_CASE_MODULE_NAME}}Model } from './{{LOWER_CASE_MODULE_NAME}}.model';
// import { {{PASCAL_CASE_MODULE_NAME}}Schema } from './{{LOWER_CASE_MODULE_NAME}}.schema';

// export class {{PASCAL_CASE_MODULE_NAME}}Service {
//   async create(data: any) {
//     const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.parse(data);
//     return await {{PASCAL_CASE_MODULE_NAME}}Model.create(validated);
//   }

//   async findAll() {
//     return await {{PASCAL_CASE_MODULE_NAME}}Model.find();
//   }

//   async findOne(id: string) {
//     const record = await {{PASCAL_CASE_MODULE_NAME}}Model.findById(id);
//     if (!record) {
//       throw new Error('{{PASCAL_CASE_MODULE_NAME}} not found');
//     }
//     return record;
//   }

//   async update(id: string, data: any) {
//     const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.partial().parse(data);
//     return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndUpdate(id, validated, { new: true });
//   }

//   async delete(id: string) {
//     return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndDelete(id);
//   }
// }
// `,
// 		},
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.model.ts",
// 			Description: "Creating Mongoose model and schema",
// 			Content: `import mongoose, { Schema, Document } from 'mongoose';

// export interface I{{PASCAL_CASE_MODULE_NAME}} extends Document {
// {{MODEL_FIELDS}}
// }

// const {{PASCAL_CASE_MODULE_NAME}}Schema = new Schema({
// {{MONGOOSE_SCHEMA_FIELDS}}
// }, { timestamps: true });

// export const {{PASCAL_CASE_MODULE_NAME}}Model = mongoose.model<I{{PASCAL_CASE_MODULE_NAME}}>('{{PASCAL_CASE_MODULE_NAME}}', {{PASCAL_CASE_MODULE_NAME}}Schema);
// `,
// 		},
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.schema.ts",
// 			Description: "Creating Zod validation schema",
// 			Content: `import { z } from 'zod';

// {{ZOD_GENERATED_SCHEMA}}

// {{ZOD_INFER_TYPES}}

// {{ZOD_EXPORTS}}
// `,
// 		},
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.routes.ts",
// 			Description: "Creating Express routes",
// 			Content: `import { Router } from 'express';
// import { {{PASCAL_CASE_MODULE_NAME}}Controller } from './{{LOWER_CASE_MODULE_NAME}}.controller';

// const router = Router();
// const controller = new {{PASCAL_CASE_MODULE_NAME}}Controller();

// router.post('/', controller.create.bind(controller));
// router.get('/', controller.findAll.bind(controller));
// router.get('/:id', controller.findOne.bind(controller));
// router.put('/:id', controller.update.bind(controller));
// router.delete('/:id', controller.delete.bind(controller));

// export default router;
// `,
// 		},
// 		{
// 			FilePath:    "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{UPPER_CASE_MODULE_NAME}}.constants.ts",
// 			Description: "Creating module constants",
// 			Content: `export const {{UPPER_CASE_MODULE_NAME}}_CONSTANTS = {
//   MODEL_NAME: '{{PASCAL_CASE_MODULE_NAME}}',
//   COLLECTION_NAME: '{{LOWER_CASE_MODULE_NAME}}s',
//   BASE_PATH: '/api/{{LOWER_CASE_MODULE_NAME}}s',
// };
// `,
// 		},
// 	}
// }

// func getFrontendInstructions() []FileInstruction {
// 	return []FileInstruction{
// 		{
// 			FilePath:    "src/features/{{LOWER_CASE_MODULE_NAME}}/components/{{PASCAL_CASE_MODULE_NAME}}List.tsx",
// 			Description: "Creating React list component",
// 			Content: `import React from 'react';
// import { use{{PASCAL_CASE_MODULE_NAME}}s } from '../hooks/use{{PASCAL_CASE_MODULE_NAME}}s';

// export const {{PASCAL_CASE_MODULE_NAME}}List: React.FC = () => {
//   const { data, isLoading, error } = use{{PASCAL_CASE_MODULE_NAME}}s();

//   if (isLoading) return <div>Loading...</div>;
//   if (error) return <div>Error: {error.message}</div>;

//   return (
//     <div>
//       <h2>{{PASCAL_CASE_MODULE_NAME}}s</h2>
//       <ul>
//         {data?.map(item => (
//           <li key={item.id}>{JSON.stringify(item)}</li>
//         ))}
//       </ul>
//     </div>
//   );
// };
// `,
// 		},
// 		{
// 			FilePath:    "src/features/{{LOWER_CASE_MODULE_NAME}}/hooks/use{{PASCAL_CASE_MODULE_NAME}}s.ts",
// 			Description: "Creating React Query hook",
// 			Content: `import { useQuery } from '@tanstack/react-query';
// import { {{LOWER_CASE_MODULE_NAME}}Api } from '../api/{{LOWER_CASE_MODULE_NAME}}.api';

// export const use{{PASCAL_CASE_MODULE_NAME}}s = () => {
//   return useQuery({
//     queryKey: ['{{LOWER_CASE_MODULE_NAME}}s'],
//     queryFn: {{LOWER_CASE_MODULE_NAME}}Api.getAll,
//   });
// };
// `,
// 		},
// 		{
// 			FilePath:    "src/features/{{LOWER_CASE_MODULE_NAME}}/api/{{LOWER_CASE_MODULE_NAME}}.api.ts",
// 			Description: "Creating API service",
// 			Content: `import axios from 'axios';

// const BASE_URL = '/api/{{LOWER_CASE_MODULE_NAME}}s';

// export const {{LOWER_CASE_MODULE_NAME}}Api = {
//   getAll: async () => {
//     const response = await axios.get(BASE_URL);
//     return response.data;
//   },

//   getOne: async (id: string) => {
//     const response = await axios.get(\${BASE_URL}/\${id});
//     return response.data;
//   },

//   create: async (data: any) => {
//     const response = await axios.post(BASE_URL, data);
//     return response.data;
//   },

//   update: async (id: string, data: any) => {
//     const response = await axios.put(\${BASE_URL}/\${id}\, data);
//     return response.data;
//   },

//   delete: async (id: string) => {
//     await axios.delete(\${BASE_URL}/\${id});
//   },
// };
// `,
// 		},
// 	}
// }

// func getUpdateInstructions(moduleType string) []UpdateInstruction {
// 	if moduleType == "bm" {
// 		return []UpdateInstruction{
// 			{
// 				FilePath:          "src/routes/index.ts",
// 				Placeholder:       "// SUPER_CLI_ROUTE_IMPORTS",
// 				Content:           "import {{LOWER_CASE_MODULE_NAME}}Routes from '../modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.routes';",
// 				Position:          "bottom",
// 				CreateIfNotExists: true,
// 				Description:       "Adding route import to main routes file",
// 			},
// 			{
// 				FilePath:          "src/routes/index.ts",
// 				Placeholder:       "// SUPER_CLI_ROUTE_DEFINITIONS",
// 				Content:           "router.use('/api/{{LOWER_CASE_MODULE_NAME}}s', {{LOWER_CASE_MODULE_NAME}}Routes);",
// 				Position:          "bottom",
// 				CreateIfNotExists: false,
// 				Description:       "Registering route in main router",
// 			},
// 		}
// 	}
// 	return []UpdateInstruction{}
// }

// func getDefaultConfig() ProjectConfig {
// 	return ProjectConfig{
// 		Name:               "super-cli-config",
// 		Version:            "1.0.0",
// 		FileInstructions:   backendmodule.GetCreateBackendModuleInstructions(),
// 		UpdateInstructions: getUpdateInstructions("bm"),
// 	}
// }

func updateFile(filePath, placeholder, content, position string, createIfNotExists bool) error {
	// Check if file exists
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		if createIfNotExists && os.IsNotExist(err) {
			// Create file with placeholder and content
			dir := filepath.Dir(filePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			newContent := fmt.Sprintf("%s\n%s\n", placeholder, content)
			return os.WriteFile(filePath, []byte(newContent), 0644)
		}
		return err
	}

	fileStr := string(fileContent)

	// Check if placeholder exists
	if !strings.Contains(fileStr, placeholder) {
		return fmt.Errorf("placeholder '%s' not found in file", placeholder)
	}

	// Check if content already exists (avoid duplicates)
	if strings.Contains(fileStr, content) {
		if verbose {
			fmt.Printf("    ⚠ Content already exists in %s, skipping\n", filePath)
		}
		return nil
	}

	// Determine insertion position
	var newContent string
	if position == "bottom" {
		// Add content after placeholder
		newContent = strings.Replace(fileStr, placeholder, placeholder+"\n"+content, 1)
	} else {
		// Default: add content before placeholder (top)
		newContent = strings.Replace(fileStr, placeholder, content+"\n"+placeholder, 1)
	}

	return os.WriteFile(filePath, []byte(newContent), 0644)
}

func buildReplacements(moduleName string, fields []types.Field) map[string]string {
	return map[string]string{
		"{{MODULE_NAME}}":                moduleName,
		"{{LOWER_CASE_MODULE_NAME}}":     strings.ToLower(moduleName),
		"{{UPPER_CASE_MODULE_NAME}}":     strings.ToUpper(moduleName),
		"{{PASCAL_CASE_MODULE_NAME}}":    toPascalCase(moduleName),
		"{{CAMEL_CASE_MODULE_NAME}}":     toCamelCase(moduleName),
		"{{MODEL_FIELDS}}":               generateModelFields(fields),
		"{{MONGOOSE_SCHEMA_FIELDS}}":     generateMongooseFields(fields),
		"{{ZOD_GENERATED_SCHEMA}}":       generateZodSchema(moduleName, fields),
		"{{ZOD_INFER_TYPES}}":            generateZodTypes(moduleName),
		"{{ZOD_EXPORTS}}":                generateZodExports(moduleName),
	}
}

func applyReplacements(content string, replacements map[string]string) string {
	result := content
	for key, value := range replacements {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

func createFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func generateModelFields(fields []types.Field) string {
	if len(fields) == 0 {
		return "  // Add your fields here"
	}
	var result strings.Builder
	for _, f := range fields {
		tsType := mapTypeToTypeScript(f.Type)
		optional := ""
		if !f.Required {
			optional = "?"
		}
		result.WriteString(fmt.Sprintf("  %s%s: %s;\n", f.Name, optional, tsType))
	}
	return strings.TrimRight(result.String(), "\n")
}

func generateMongooseFields(fields []types.Field) string {
	if len(fields) == 0 {
		return "  // Add your fields here"
	}
	var result strings.Builder
	for _, f := range fields {
		mongoType := mapTypeToMongoose(f.Type)
		result.WriteString(fmt.Sprintf("  %s: { type: %s, required: %t },\n", f.Name, mongoType, f.Required))
	}
	return strings.TrimRight(result.String(), "\n")
}

func generateZodSchema(moduleName string, fields []types.Field) string {
	if len(fields) == 0 {
		return "export const " + toPascalCase(moduleName) + "Schema = z.object({\n  // Add your fields here\n});"
	}
	var result strings.Builder
	result.WriteString("export const " + toPascalCase(moduleName) + "Schema = z.object({\n")
	for _, f := range fields {
		zodType := mapTypeToZod(f.Type)
		if !f.Required {
			zodType += ".optional()"
		}
		result.WriteString(fmt.Sprintf("  %s: %s,\n", f.Name, zodType))
	}
	result.WriteString("});")
	return result.String()
}

func generateZodTypes(moduleName string) string {
	return fmt.Sprintf("export type %s = z.infer<typeof %sSchema>;", toPascalCase(moduleName), toPascalCase(moduleName))
}

func generateZodExports(moduleName string) string {
	return fmt.Sprintf("export { %sSchema };", toPascalCase(moduleName))
}

func mapTypeToTypeScript(t string) string {
	switch t {
	case "S":
		return "string"
	case "N":
		return "number"
	case "B":
		return "boolean"
	case "D":
		return "Date"
	default:
		return "any"
	}
}

func mapTypeToMongoose(t string) string {
	switch t {
	case "S":
		return "String"
	case "N":
		return "Number"
	case "B":
		return "Boolean"
	case "D":
		return "Date"
	default:
		return "Schema.Types.Mixed"
	}
}

func mapTypeToZod(t string) string {
	switch t {
	case "S":
		return "z.string()"
	case "N":
		return "z.number()"
	case "B":
		return "z.boolean()"
	case "D":
		return "z.date()"
	default:
		return "z.any()"
	}
}