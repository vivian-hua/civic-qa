import React, { useEffect, useState } from 'react';
import { Header } from '../Components/Header';
import { SubDashboard, InquiryData } from '../Components/SubDashboard';
import { SubHeaderLine } from '../Components/SubHeaderLine';
import { StatCardRow } from '../Components/StatCardRow';
import { Inquiries } from './Inquiries';
import * as Endpoints from '../Constants/Endpoints';

export function Casework() {
    const auth = localStorage.getItem("Authorization") || "";
    const [onInquiriesView, setInquiriesView] = useState(false);
    const [specificViewTitle, setSpecificViewTitle] = useState("");
    const [specificTopicsData, setSpecificTopicsData] = useState<InquiryData[]>([]);
    const [topicsData, setTopicsData] = useState<Map<string, InquiryData[]>>();
    const [topicsCases, setTopicsCases] = useState<InquiryData[]>([]);
    const [summaryToday, setSummaryToday] = useState(0);
    const [summaryTotal, setSummaryTotal] = useState(0);
    const [summaryTopics, setSummaryTopics] = useState(0);


    const getInquiries = async() => {
        const response = await fetch(Endpoints.Base + Endpoints.ResponsesActiveCasework, {
            method: "GET",
            headers: new Headers({
                "Authorization": auth
            })
        });
        if (response.status >= 300) {
            console.log("Error retrieving form responses");
            return;
        }
        const inquiriesCasework = await response.json();
        var inquiries: InquiryData[] = [];
        let topicsMap = new Map<string, InquiryData[]>();
        let topicsCases = new Map<string, number>();

        inquiriesCasework.forEach(function(inquiry: any) {
            var d = new Date(inquiry.createdAt);
            var t = d.toLocaleString("en-US");
            var topics = inquiry.tags;
            var data: InquiryData = {id: inquiry.id, email: inquiry.emailAddress, name: inquiry.name + " / " + inquiry.subject, value: t, body: inquiry.body}

            topics.forEach((topic: any) => {
                if (topicsMap.has(topic.value)) {
                    var getList = topicsMap.get(topic.value);
                    getList?.push(data);
                    topicsMap.set(topic.value, getList || []);
                } else {
                    var newList: InquiryData[] = [];
                    newList.push(data);
                    topicsMap.set(topic.value, newList);
                }

                topicsCases.set(topic.value, (topicsCases.get(topic.value) || 0) + 1);

            });
            inquiries.push(data);
        });

        var cases: InquiryData[] = [];
        Array.from(topicsCases.keys()).forEach((key) => {
            var subText = " case";
            if ((topicsCases.get(key) || 0) > 1) {
                subText = " cases";
            }
            cases.push({name: key, value: topicsCases.get(key) + subText});
        })

        cases.sort((a, b) => (a.value > b.value) ? -1 : (a.value === b.value) ? -1 : 1);

        setSummaryTopics(cases.length);
        setSummaryTotal(inquiries.length);
        setTopicsCases(cases);
        setTopicsData(topicsMap);
    }

    const getInquiriesToday = async() => {
        const response = await fetch(Endpoints.Base + Endpoints.ResponsesActiveCaseworkTodayOnly, {
            method: "GET",
            headers: new Headers({
                "Authorization": auth
            })
        });
        if (response.status >= 300) {
            console.log("Error retrieving form responses");
            return;
        }
        const inquiriesToday = await response.json();
        setSummaryToday(inquiriesToday.length);
    }

    function inquiriesView(data: InquiryData) {
        setSpecificViewTitle(data.name);
        setSpecificTopicsData(topicsData?.get(data.name) || []);
        setInquiriesView(true);
    }

    function initialView() {
        setInquiriesView(false);
        getInquiries();
    }
    
    useEffect(() => {
        getInquiries();
        getInquiriesToday();
    }, []);

    let statCards = [
        {title: "New Today", stat: summaryToday},
        {title: "Total", stat: summaryTotal},
        {title: "Topics", stat: summaryTopics}
    ];

    return (
        onInquiriesView ? 
            <div> 
                <div className="dashboard sub-dashboard">
                    <button className="exit-button" onClick={initialView}><img className="back-arrow" src="./assets/icons/arrow.svg"></img></button>
                </div>
                <Inquiries header="Casework" subjectTitle={specificViewTitle} data={specificTopicsData} hideInquiryBackArrow={true}></Inquiries>
            </div>
            : <div className="dashboard sub-dashboard">
                <div>
                    <Header title="Casework Topics"></Header>
                    <SubDashboard title="TOPICS" data={topicsCases} changeViewFunc={inquiriesView} fullPageView={false}></SubDashboard>
                    <div className="sub-summary">
                        <SubHeaderLine title="SUMMARY" subHeaderValue={"Active Cases"}></SubHeaderLine>
                        <StatCardRow spaceEven={false} cards={statCards}></StatCardRow>
                    </div>
                </div>
            </div>
    );
}