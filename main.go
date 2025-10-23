package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uploadCmd)
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

	fmt.Print("Enter module name: ")
	fmt.Scanln(&moduleName)

	if moduleName == "" {
		fmt.Println("Error: module name is required")
		return
	}

	fmt.Print("Enter fields (optional, format: name@S@R,email@S@R): ")
	fmt.Scanln(&fields)

	module := Module{
		ModuleName:      moduleName,
		ModelProperties: fields,
	}

	generateModule(moduleType, module)
	fmt.Printf("✓ Module '%s' created successfully!\n", moduleName)
}

func handleBulkUpload(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var modules []Module
	if err := json.Unmarshal(data, &modules); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Creating %d modules...\n", len(modules))
	for _, module := range modules {
		generateModule("bm", module)
		fmt.Printf("✓ Module '%s' created\n", module.ModuleName)
	}
	fmt.Println("\n✓ Bulk creation completed!")
}

type Module struct {
	ModuleName      string `json:"moduleName"`
	ModelProperties string `json:"modelProperties"`
}

type Field struct {
	Name     string
	Type     string
	Required bool
}

func parseFields(fieldsStr string) []Field {
	if fieldsStr == "" {
		return []Field{}
	}

	fieldParts := strings.Split(fieldsStr, ",")
	fields := make([]Field, 0, len(fieldParts))

	for _, part := range fieldParts {
		segments := strings.Split(strings.TrimSpace(part), "@")
		if len(segments) >= 2 {
			field := Field{
				Name:     segments[0],
				Type:     segments[1],
				Required: len(segments) >= 3 && segments[2] == "R",
			}
			fields = append(fields, field)
		}
	}

	return fields
}

func generateModule(moduleType string, module Module) {
	instructions := getInstructions(moduleType)
	fields := parseFields(module.ModelProperties)

	replacements := buildReplacements(module.ModuleName, fields)

	for _, instruction := range instructions {
		filePath := applyReplacements(instruction.FilePath, replacements)
		fileContent := applyReplacements(instruction.Content, replacements)

		if err := createFile(filePath, fileContent); err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			continue
		}
		fmt.Printf("  ✓ Created: %s\n", filePath)
	}
}

type FileInstruction struct {
	FilePath string
	Content  string
}

func getInstructions(moduleType string) []FileInstruction {
	if moduleType == "bm" {
		return getBackendInstructions()
	}
	return getFrontendInstructions()
}

func getBackendInstructions() []FileInstruction {
	return []FileInstruction{
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.controller.ts",
			Content: `import { Request, Response } from 'express';
import { {{PASCAL_CASE_MODULE_NAME}}Service } from './{{LOWER_CASE_MODULE_NAME}}.service';

export class {{PASCAL_CASE_MODULE_NAME}}Controller {
  private service: {{PASCAL_CASE_MODULE_NAME}}Service;

  constructor() {
    this.service = new {{PASCAL_CASE_MODULE_NAME}}Service();
  }

  async create(req: Request, res: Response) {
    try {
      const data = await this.service.create(req.body);
      res.status(201).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async findAll(req: Request, res: Response) {
    try {
      const data = await this.service.findAll();
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async findOne(req: Request, res: Response) {
    try {
      const data = await this.service.findOne(req.params.id);
      res.status(200).json(data);
    } catch (error) {
      res.status(404).json({ error: error.message });
    }
  }

  async update(req: Request, res: Response) {
    try {
      const data = await this.service.update(req.params.id, req.body);
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }

  async delete(req: Request, res: Response) {
    try {
      await this.service.delete(req.params.id);
      res.status(204).send();
    } catch (error) {
      res.status(500).json({ error: error.message });
    }
  }
}
`,
		},
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.service.ts",
			Content: `import { {{PASCAL_CASE_MODULE_NAME}}Model } from './{{LOWER_CASE_MODULE_NAME}}.model';
import { {{PASCAL_CASE_MODULE_NAME}}Schema } from './{{LOWER_CASE_MODULE_NAME}}.schema';

export class {{PASCAL_CASE_MODULE_NAME}}Service {
  async create(data: any) {
    const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.parse(data);
    return await {{PASCAL_CASE_MODULE_NAME}}Model.create(validated);
  }

  async findAll() {
    return await {{PASCAL_CASE_MODULE_NAME}}Model.find();
  }

  async findOne(id: string) {
    const record = await {{PASCAL_CASE_MODULE_NAME}}Model.findById(id);
    if (!record) {
      throw new Error('{{PASCAL_CASE_MODULE_NAME}} not found');
    }
    return record;
  }

  async update(id: string, data: any) {
    const validated = {{PASCAL_CASE_MODULE_NAME}}Schema.partial().parse(data);
    return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndUpdate(id, validated, { new: true });
  }

  async delete(id: string) {
    return await {{PASCAL_CASE_MODULE_NAME}}Model.findByIdAndDelete(id);
  }
}
`,
		},
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.model.ts",
			Content: `import mongoose, { Schema, Document } from 'mongoose';

export interface I{{PASCAL_CASE_MODULE_NAME}} extends Document {
{{MODEL_FIELDS}}
}

const {{PASCAL_CASE_MODULE_NAME}}Schema = new Schema({
{{MONGOOSE_SCHEMA_FIELDS}}
}, { timestamps: true });

export const {{PASCAL_CASE_MODULE_NAME}}Model = mongoose.model<I{{PASCAL_CASE_MODULE_NAME}}>('{{PASCAL_CASE_MODULE_NAME}}', {{PASCAL_CASE_MODULE_NAME}}Schema);
`,
		},
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.schema.ts",
			Content: `import { z } from 'zod';

{{ZOD_GENERATED_SCHEMA}}

{{ZOD_INFER_TYPES}}

{{ZOD_EXPORTS}}
`,
		},
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{LOWER_CASE_MODULE_NAME}}.routes.ts",
			Content: `import { Router } from 'express';
import { {{PASCAL_CASE_MODULE_NAME}}Controller } from './{{LOWER_CASE_MODULE_NAME}}.controller';

const router = Router();
const controller = new {{PASCAL_CASE_MODULE_NAME}}Controller();

router.post('/', controller.create.bind(controller));
router.get('/', controller.findAll.bind(controller));
router.get('/:id', controller.findOne.bind(controller));
router.put('/:id', controller.update.bind(controller));
router.delete('/:id', controller.delete.bind(controller));

export default router;
`,
		},
		{
			FilePath: "src/modules/{{LOWER_CASE_MODULE_NAME}}/{{UPPER_CASE_MODULE_NAME}}.constants.ts",
			Content: `export const {{UPPER_CASE_MODULE_NAME}}_CONSTANTS = {
  MODEL_NAME: '{{PASCAL_CASE_MODULE_NAME}}',
  COLLECTION_NAME: '{{LOWER_CASE_MODULE_NAME}}s',
  BASE_PATH: '/api/{{LOWER_CASE_MODULE_NAME}}s',
};
`,
		},
	}
}

func getFrontendInstructions() []FileInstruction {
	return []FileInstruction{
		{
			FilePath: "src/features/{{LOWER_CASE_MODULE_NAME}}/components/{{PASCAL_CASE_MODULE_NAME}}List.tsx",
			Content: `import React from 'react';
import { use{{PASCAL_CASE_MODULE_NAME}}s } from '../hooks/use{{PASCAL_CASE_MODULE_NAME}}s';

export const {{PASCAL_CASE_MODULE_NAME}}List: React.FC = () => {
  const { data, isLoading, error } = use{{PASCAL_CASE_MODULE_NAME}}s();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>{{PASCAL_CASE_MODULE_NAME}}s</h2>
      <ul>
        {data?.map(item => (
          <li key={item.id}>{JSON.stringify(item)}</li>
        ))}
      </ul>
    </div>
  );
};
`,
		},
		{
			FilePath: "src/features/{{LOWER_CASE_MODULE_NAME}}/hooks/use{{PASCAL_CASE_MODULE_NAME}}s.ts",
			Content: `import { useQuery } from '@tanstack/react-query';
import { {{LOWER_CASE_MODULE_NAME}}Api } from '../api/{{LOWER_CASE_MODULE_NAME}}.api';

export const use{{PASCAL_CASE_MODULE_NAME}}s = () => {
  return useQuery({
    queryKey: ['{{LOWER_CASE_MODULE_NAME}}s'],
    queryFn: {{LOWER_CASE_MODULE_NAME}}Api.getAll,
  });
};
`,
		},
		{
			FilePath: "src/features/{{LOWER_CASE_MODULE_NAME}}/api/{{LOWER_CASE_MODULE_NAME}}.api.ts",
			Content: `import axios from 'axios';

const BASE_URL = '/api/{{LOWER_CASE_MODULE_NAME}}s';

export const {{LOWER_CASE_MODULE_NAME}}Api = {
  getAll: async () => {
    const response = await axios.get(BASE_URL);
    return response.data;
  },

  getOne: async (id: string) => {
    const response = await axios.get(\'\${BASE_URL}/\${id}\');
    return response.data;
  },

  create: async (data: any) => {
    const response = await axios.post(BASE_URL, data);
    return response.data;
  },

  update: async (id: string, data: any) => {
    const response = await axios.put(\'\${BASE_URL}/\${id}\', data);
    return response.data;
  },

  delete: async (id: string) => {
    await axios.delete("${BASE_URL}/${id}");
  },
};
`,
		},
	}
}

func buildReplacements(moduleName string, fields []Field) map[string]string {
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
	dir := path[:strings.LastIndex(path, "/")]
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

func generateModelFields(fields []Field) string {
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

func generateMongooseFields(fields []Field) string {
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

func generateZodSchema(moduleName string, fields []Field) string {
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