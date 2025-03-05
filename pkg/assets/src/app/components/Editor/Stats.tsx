import React, { useState, useEffect } from "react";
import {
  LineChart,
  Line,
  YAxis,
  ResponsiveContainer,
  PieChart,
  Pie,
  Tooltip,
  Cell,
  Legend,
} from "recharts";
import { XIcon, RefreshCwIcon } from "lucide-react";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { DocStats, OpStats } from "@/../rogueEditor";
import { Spinner } from "../ui/spinner";
import { Button } from "../ui/button";

interface StatsProps {
  onClose?: () => void;
}

const generateChartData = (dataArray: number[]) =>
  dataArray.map((value, index) => ({ index, value }));

export const Stats: React.FC<StatsProps> = function ({ onClose }) {
  const { editor, editorMode } = useRogueEditorContext();
  const [opStats, setOpStats] = useState<OpStats | null>(null);
  const [docStats, setDocStats] = useState<DocStats | null>(null);

  const fetchStats = async () => {
    if (editor) {
      // Fetch new stats from the editor
      const newStats = await editor.getOpStats();
      setOpStats(newStats);
      const newDocStats = await editor.getDocStats();
      setDocStats(newDocStats);
    }
  };

  useEffect(() => {
    if (!editor || editorMode !== "xray") {
      return;
    }
    fetchStats();
  }, [editor, editorMode]);

  const handleReload = () => {
    setOpStats(null);
    setDocStats(null);
    fetchStats();
  };

  if (!editor || editorMode !== "xray") {
    return null;
  }

  if (!opStats || !docStats) {
    return (
      <div className="w-full flex items-center justify-center">
        <Spinner />
      </div>
    );
  }

  const {
    inserts,
    deletes,
    insertsByPrefix,
    deletesByPrefix,
    currentCharsByPrefix,
  } = opStats;

  const { wordCount, paragraphCount } = docStats;

  // Prefix to Title Mapping
  const prefixTitles: { [key: string]: string } = {
    "": "User",
    "!": "AI",
    $: "Paste",
    "#": "External Paste",
  };

  // Define prefix order
  const prefixOrder = ["", "!", "$", "#"];

  // Collect all unique prefixes from insertsByPrefix and deletesByPrefix
  const dataPrefixes = new Set([
    ...Object.keys(insertsByPrefix),
    ...Object.keys(deletesByPrefix),
  ]);

  // Get prefixes to display, in the desired order
  const prefixesToDisplay = prefixOrder.filter((prefix) =>
    dataPrefixes.has(prefix),
  );

  // Collect all data arrays for min and max computation
  const allDataArrays = [
    inserts,
    deletes,
    ...Object.values(insertsByPrefix),
    ...Object.values(deletesByPrefix),
  ];

  // Flatten all values into a single array
  const allValues = allDataArrays.flat();

  // Compute global min and max values
  const minY = allValues.length > 0 ? Math.min(...allValues) : 0;
  const maxY = allValues.length > 0 ? Math.max(...allValues) : 1;

  const sparklineProps = {
    width: "100%",
    height: 50, // Adjust height for a compact sparkline
  };

  const lineProps = {
    type: "monotone",
    dataKey: "value",
    strokeWidth: 1,
    dot: false, // Remove dots from the line
  };

  // Prepare pie chart data
  const pieData = Object.entries(currentCharsByPrefix).map(
    ([prefix, value]) => ({
      prefix,
      name: prefixTitles[prefix] || `Prefix "${prefix}"`,
      value,
    }),
  );

  // Define colors for each prefix
  const prefixColors: { [key: string]: string } = {
    "": "hsla(201, 95%, 52%, 0.5)", // User
    "!": "hsla(258.3, 90%, 66%, 0.5)", // AI
    $: "hsla(163, 94%, 37%, 0.5)", // Paste
    "#": "hsla(163, 94%, 37%, 0.5)", // External Paste
  };

  return (
    <div className="flex flex-col space-y-4 pointer-events-auto bg-background rounded rounded-lg p-4">
      {/* Buttons */}
      <div className="flex justify-center space-x-2">
        <div className="text-lg font-bold mt-1">Word Count: {wordCount}</div>
        <div className="flex-grow"></div>
        <Button variant="ghost" onClick={handleReload}>
          <RefreshCwIcon className="w-4 h-4" />
        </Button>
        {onClose && (
          <Button variant="ghost" onClick={onClose}>
            <XIcon className="w-4 h-4" />
          </Button>
        )}
      </div>

      <div className="flex flex-row space-x-4">
        <div className="flex flex-col space-y-4 w-full">
          {/* Inserts and Deletes Charts without Prefix */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3>Total Inserts</h3>
              <ResponsiveContainer {...sparklineProps}>
                <LineChart
                  data={generateChartData(inserts)}
                  margin={{ top: 0, right: 0, bottom: 0, left: 0 }}
                >
                  <YAxis domain={[minY, maxY]} hide={true} />
                  <Line {...lineProps} stroke="#777" />
                </LineChart>
              </ResponsiveContainer>
            </div>

            <div>
              <h3>Total Deletes</h3>
              <ResponsiveContainer {...sparklineProps}>
                <LineChart
                  data={generateChartData(deletes)}
                  margin={{ top: 0, right: 0, bottom: 0, left: 0 }}
                >
                  <YAxis domain={[minY, maxY]} hide={true} />
                  <Line {...lineProps} stroke="#777" />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Grouped Inserts and Deletes by Prefix */}
          {prefixesToDisplay.map((prefix) => {
            const insertsData = insertsByPrefix[prefix] || [];
            const deletesData = deletesByPrefix[prefix] || [];
            const prefixTitle = prefixTitles[prefix] || `Prefix "${prefix}"`;

            return (
              <div key={`group-${prefix}`} className="grid grid-cols-2 gap-4">
                <div>
                  <h3>{prefixTitle} Inserts</h3>
                  <ResponsiveContainer {...sparklineProps}>
                    <LineChart
                      data={generateChartData(insertsData)}
                      margin={{ top: 0, right: 0, bottom: 0, left: 0 }}
                    >
                      <YAxis domain={[minY, maxY]} hide={true} />
                      <Line {...lineProps} stroke={prefixColors[prefix]} />
                    </LineChart>
                  </ResponsiveContainer>
                </div>

                <div>
                  <h3>{prefixTitle} Deletes</h3>
                  <ResponsiveContainer {...sparklineProps}>
                    <LineChart
                      data={generateChartData(deletesData)}
                      margin={{ top: 0, right: 0, bottom: 0, left: 0 }}
                    >
                      <YAxis domain={[minY, maxY]} hide={true} />
                      <Line {...lineProps} stroke={prefixColors[prefix]} />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </div>
            );
          })}
        </div>

        {/* Pie Chart */}
        <div className="flex flex-col items-center">
          <ResponsiveContainer width={300} height={350}>
            <PieChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
              <Pie
                data={pieData}
                dataKey="value"
                nameKey="name"
                cx="50%"
                cy="45%"
                outerRadius={100}
                fill="#8884d8"
                labelLine={false}
                label={({ percent }) => `${(percent * 100).toFixed(0)}%`}
              >
                {pieData.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={prefixColors[entry.prefix] || "#8884d8"}
                  />
                ))}
              </Pie>
              <Tooltip />
              <Legend
                layout="horizontal"
                verticalAlign="bottom"
                align="center"
              />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};
