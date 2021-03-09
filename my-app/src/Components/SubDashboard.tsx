import { Dispatch, SetStateAction } from 'react';
import { SubDashboardCard } from './SubDashboardCard';
import { SubHeaderLine } from './SubHeaderLine';
import './SubDashboard.css';
import { EmailTemplateCard } from './EmailTemplateCard';

export type SubDashboardData = {
    id: string;
    name: string;
    value: any;
    body?: string;
}

export type SubDashboardProps = {
    title: string;
    data: Array<SubDashboardData>;
    changeViewFunc: Function;
    emailTemplates: boolean;
    fullPageView: boolean;
    subHeaderNumber?: number;
};

export function SubDashboard(props: SubDashboardProps) {
    let cards:any[] = [];
    if (props.emailTemplates) {
        props.data.forEach(d => cards.push(<EmailTemplateCard name={d.name} value={d.value}></EmailTemplateCard>))
    } else {
        props.data.forEach(d => cards.push(<SubDashboardCard data={d} changeViewFunc={props.changeViewFunc}></SubDashboardCard>));
    }

    return (
        <div>
            <div>
                <SubHeaderLine title={props.title} subHeaderNumber={props.subHeaderNumber ? props.subHeaderNumber : undefined}></SubHeaderLine>
                {props.fullPageView ? 
                <div className="sub-dash-cards-700">
                    {cards}
                </div> :
                <div className="sub-dash-cards-400">
                    {cards}
                </div>
                }

            </div>
        </div>
    );
}