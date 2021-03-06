import { useEffect, useState } from 'react';
import { Header } from '../Components/Header';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import * as Endpoints from '../Constants/Endpoints';

export type EngagementReportBarChartData = {
    name: string;
    count: number;
}

export function EngagementReports() {
  const [engagementChartData, setEngagementChartData] = useState<EngagementReportBarChartData[]>();

    const getTags = async() => {
      var authToken = localStorage.getItem("Authorization") || "";
      const response = await fetch(Endpoints.Base + Endpoints.Tags, {
          method: "GET",
          headers: new Headers({
              "Authorization": authToken
          })
      });
      if (response.status >= 300) {
          console.log("Error retrieving form responses");
          return;
      }
      var topicCountsMap = new Map<string, number>();
      const tags = await response.json();
      tags.forEach((tag: any) => {
        if (topicCountsMap.has(tag.value)) {
          topicCountsMap.set(tag.value, (topicCountsMap.get(tag.value) || 0) + 1);
        } else {
          topicCountsMap.set(tag.value, 1);
        }
      });
      
      var engagementChartData: EngagementReportBarChartData[] = [];

      Array.from(topicCountsMap.keys()).forEach((key) => {
        engagementChartData.push({name: key, count: topicCountsMap.get(key) || 0});
      });

      engagementChartData.sort((a, b) => {
        if (a.count > b.count) {
            return -1;
        } else if (a.count < b.count) {
            return 1;
        } else {
            return 0;
        }
      });

      engagementChartData = engagementChartData.slice(0, 5);
      setEngagementChartData(engagementChartData);
    }

    useEffect(() => {
        getTags();
    }, []);

    return (
        <div className="dashboard sub-dashboard view">
            <Header title="Engagement Reports"></Header>
            <h2>Most Frequent Topics</h2>
            <div style={{ width: '100%', height: 500 }}>
              <ResponsiveContainer>
                <BarChart
                    data={engagementChartData}
                    margin={{
                        top: 5,
                        right: 30,
                        left: 20,
                        bottom: 5,
                    }}
                    >
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="count" fill="#855CF8" />
                </BarChart>
              </ResponsiveContainer>
            </div>
        </div>
    );
}